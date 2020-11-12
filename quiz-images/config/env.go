package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// GetEnvVariable from the env file
func GetEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
