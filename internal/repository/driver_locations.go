package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/mrtuuro/driver-location-api/internal/apperror"
	"github.com/mrtuuro/driver-location-api/internal/code"
	"github.com/mrtuuro/driver-location-api/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DriverRepository interface {
	CreateMany(ctx context.Context, drivers []model.DriverLocation) error
	FindNearbyDrivers(ctx context.Context,
		point *model.GeoPoint,
		radiusMeters float64, limit int) ([]model.DriverWithDistance, error)
}

type mongoDriverRepository struct {
	collection *mongo.Collection
}

func NewMongoDriverRepository(coll *mongo.Collection) DriverRepository {
	ensureIndexes(coll)
	return &mongoDriverRepository{collection: coll}
}

func ensureIndexes(col *mongo.Collection) {
	mod := mongo.IndexModel{
		Keys: bson.D{{Key: "location", Value: "2dsphere"}},
	}
	if _, err := col.Indexes().CreateOne(context.Background(), mod); err != nil {
		log.Println("Failed to create 2dsphere index:", err)
	}
}

func (r *mongoDriverRepository) CreateMany(ctx context.Context, drivers []model.DriverLocation) error {
	_, err := r.collection.InsertMany(ctx, drivers)
	if err != nil {
		return apperror.NewAppError(
			code.ErrSystemDBFailure,
			err,
			code.GetErrorMessage(code.ErrSystemDBFailure),
			)
	}
	return nil
}

func (r *mongoDriverRepository) FindNearbyDrivers(ctx context.Context, point *model.GeoPoint, radiusMeters float64, limit int) ([]model.DriverWithDistance, error) {

	pipeline := mongo.Pipeline{
		{{"$geoNear", bson.D{
			{"near", point},
			{"distanceField", "distanceMeters"},
			{"maxDistance", radiusMeters},
			{"spherical", true},
		}}},
		{{"$limit", limit}},
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, apperror.NewAppError(
			code.ErrSystemDBFailure,
			err,
			code.GetErrorMessage(code.ErrSystemDBFailure),
			)
	}
	defer cur.Close(ctx)

	var driverLocs []model.DriverWithDistance
	if err := cur.All(ctx, &driverLocs); err != nil {
		fmt.Println("here mongo problem")
		return nil, apperror.NewAppError(
			code.ErrSystemInternal,
			err,
			code.GetErrorMessage(code.ErrSystemInternal),
			)
	}
	return driverLocs, nil
}
