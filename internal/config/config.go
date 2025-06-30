package config

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	MongoURI       string
	DatabaseName   string
	CollectionName string
	Ctx            context.Context
	SecretKey      string
}

func NewConfig() *Config {
	if err := loadEnv(); err != nil {
		log.Printf("Error loading env: %v", err)
		return nil
	}

	ctx := context.Background()
	cfg := &Config{
		Port:           getEnvWithDefault("PORT", "10001"),
		MongoURI:       getEnvWithDefault("MONGO_URI", ""),
		DatabaseName:   getEnvWithDefault("DATABASE_NAME", "driver-matcher"),
		CollectionName: getEnvWithDefault("COLLECTION_NAME", "driver-locaations"),
		Ctx:            ctx,
		SecretKey:      getEnvWithDefault("SECRET_KEY", "driver-rider-matching-api-secret-key"),
	}

	return cfg
}

func getEnvWithDefault(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func loadEnv() error {
	return godotenv.Load()
}
