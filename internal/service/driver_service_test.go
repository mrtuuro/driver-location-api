package service

import (
	"context"
	"testing"

	"github.com/mrtuuro/driver-location-api/internal/model"
	"github.com/mrtuuro/driver-location-api/internal/repository"
	"github.com/mrtuuro/driver-location-api/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) CreateMany(ctx context.Context, drivers []model.DriverLocation) error {
	return m.Called(ctx, drivers).Error(0)
}
func (m *mockRepo) FindNearbyDrivers(ctx context.Context, point *model.GeoPoint, radiusMeters float64, limit int) ([]model.DriverWithDistance, error) {
	args := m.Called(ctx, point, radiusMeters, limit)
	return args.Get(0).([]model.DriverWithDistance), args.Error(1)
}

func TestSearchNearby_NoResults(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	svc := NewDriverService(repo)

	payload := model.GeoPoint{Type: "Point", Coordinates: [2]float64{41, 29}}
	repo.On("FindNearbyDrivers", ctx, &payload, 2000.0, 10).Return([]model.DriverWithDistance{}, nil)

	res, err := svc.SearchDriver(ctx, &payload, 2000.0, 10)
	assert.Error(t, err)
	assert.Empty(t, res)
	repo.AssertExpectations(t)
}

func TestAddDrivers_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	svc := NewDriverService(repo)

	payload := []model.DriverLocation{
		{DriverID: "d1", Location: model.GeoPoint{Type: "Point", Coordinates: [2]float64{41.209, 29.093}}},
		{DriverID: "d2", Location: model.GeoPoint{Type: "Point", Coordinates: [2]float64{41.208, 29.093}}},
		{DriverID: "d3", Location: model.GeoPoint{Type: "Point", Coordinates: [2]float64{41.209, 29.0919}}},
	}
	repo.On("CreateMany", ctx, payload).Return(nil)

	err := svc.AddDrivers(ctx, payload)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestSearchNearbySuccess(t *testing.T) {
	ctx := context.Background()
	repo := new(mockRepo)
	svc := NewDriverService(repo)

	userPt := &model.GeoPoint{Type: "Point", Coordinates: [2]float64{41, 29}}
	want := []model.DriverWithDistance{
		{DriverLocation: model.DriverLocation{DriverID: "d1"}, DistanceMeters: 87.0},
	}

	repo.On("FindNearbyDrivers", ctx, userPt, 2000.0, 1).Return(want, nil)

	got, err := svc.SearchDriver(ctx, userPt, 2000.0, 1)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
	repo.AssertExpectations(t)
}

func TestFindNearbyDriver_CorrectOrder(t *testing.T) {
	ctx := context.Background()

	uri, kill := test.StartMongo(ctx)
	defer kill()

	cli, _ := mongo.Connect(options.Client().ApplyURI(uri))
	col := cli.Database("geo").Collection("drivers")
	repo := repository.NewMongoDriverRepository(col)

	drivers := []model.DriverLocation{
		{
			DriverID: "A",
			Location: model.GeoPoint{Type: "Point", Coordinates: [2]float64{13.4050, 52.5200}},
		},
		{
			DriverID: "B",
			Location: model.GeoPoint{Type: "Point", Coordinates: [2]float64{13.7331, 52.4147}},
		},
	}
	require.NoError(t, repo.CreateMany(ctx, drivers))

	user := &model.GeoPoint{Type: "Point", Coordinates: [2]float64{13.3777, 52.5163}}
	list, err := repo.FindNearbyDrivers(ctx, user, 30000, 2) 
	require.NoError(t, err)
	require.Len(t, list, 2)

	require.Equal(t, "A", list[0].DriverID)
	require.Greater(t, list[1].DistanceMeters, list[0].DistanceMeters)
}
