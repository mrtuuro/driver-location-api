package model

import "github.com/google/uuid"

type DriverLocation struct {
	DriverID string   `bson:"driverId" json:"driverId"`
	Location GeoPoint `bson:"location" json:"location"`
}

type GeoPoint struct {
	Type        string     `bson:"type" json:"type"`               // Point
	Coordinates [2]float64 `bson:"coordinates" json:"coordinates"` // [latitude, longitude] but MongoDB accepts [longitude, latitude] if sending individually
}

func NewDriverLocation(lat, long float64) *DriverLocation {
	dl := &DriverLocation{}
	dl.DriverID = uuid.NewString()
	dl.Location.Type = "Point"
	dl.Location.Coordinates[0] = lat
	dl.Location.Coordinates[1] = long
	return dl
}

type DriverWithDistance struct {
	DriverLocation `bson:",inline" json:",inline"`
	DistanceMeters float64 `bson:"distanceMeters" json:"distanceMeters"`
}
