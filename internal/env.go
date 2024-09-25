package internal

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	err := godotenv.Load()
	if err != nil && os.Getenv("MODE") == "development" {
		log.Fatal("Error reading .env file!")
	}

	log.Printf("env-%s: %s", key, os.Getenv(key))

	return strings.TrimSpace(os.Getenv(key))
}
