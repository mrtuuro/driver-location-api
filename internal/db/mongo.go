package db

import (
	"github.com/mrtuuro/driver-location-api/internal/apperror"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(uri string) (*mongo.Client, error) {
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, apperror.NewAppError(
			"ERR_CONNECT_CLIENT",
			err,
			"An error occurred connecting mongodb.",
			)
	}
	return client, nil
}

func GetCollection(client *mongo.Client, dbName, collName string) *mongo.Collection {
	db := client.Database(dbName)
	coll := db.Collection(collName, nil)
	return coll
}
