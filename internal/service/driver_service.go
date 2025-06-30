package service

import (
	"context"
	"errors"

	"github.com/mrtuuro/driver-location-api/internal/apperror"
	"github.com/mrtuuro/driver-location-api/internal/code"
	"github.com/mrtuuro/driver-location-api/internal/model"
	"github.com/mrtuuro/driver-location-api/internal/repository"
)

type DriverService interface {
	AddDrivers(ctx context.Context, drivers []model.DriverLocation) error
	SearchDriver(ctx context.Context, geoData *model.GeoPoint, radius float64, limit int) ([]model.DriverWithDistance, error)
}

type driverService struct {
	repo repository.DriverRepository
}

func NewDriverService(repo repository.DriverRepository) DriverService {
	return &driverService{repo: repo}
}

func (svc *driverService) AddDrivers(ctx context.Context, drivers []model.DriverLocation) error {
	return svc.repo.CreateMany(ctx, drivers)
}

func (svc *driverService) SearchDriver(ctx context.Context, geoData *model.GeoPoint, radius float64, limit int) ([]model.DriverWithDistance, error) {
	drivers, err := svc.repo.FindNearbyDrivers(ctx, geoData, radius, limit)
	if err != nil {
		return nil, err
	}

	if len(drivers) == 0 {
		return nil, apperror.NewAppError(
			code.ErrNearbyDriverNotFound,
			errors.New("Db operation success but driver not found"),
			code.GetErrorMessage(code.ErrNearbyDriverNotFound),
			)
	}
	return drivers, nil
}
