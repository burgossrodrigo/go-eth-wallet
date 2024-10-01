package config

import (
	"log"
	"os"

	models "wallet/pkg/models"

	"github.com/joho/godotenv"
)

func LoadEnv() (cfg models.Config) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	cfg.MongoURI = os.Getenv("mongoURI")
	cfg.JWTSecret = os.Getenv("jwtSecret")

	return cfg
}
