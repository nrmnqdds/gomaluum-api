package internal

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	err := godotenv.Load()
	if err != nil && os.Getenv("ENV") == "development" {
		log.Fatal("Error reading .env file!")
	}

	log.Printf("env-%s: %s", key, os.Getenv(key))

	return os.Getenv(key)
}
