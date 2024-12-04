package services

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/helpers"
)

func CatalogScraper() (interface{}, *dtos.CustomError) {
	logger, _ := helpers.NewLogger()

	basepath, err := os.Getwd()
	if err != nil {
		logger.Error(err)
	}

	logger.Info(basepath)

	path := filepath.Join(basepath, "dtos/iium_2024_2025_1.json")

	if os.Getenv("ENVIRONMENT") == "production" {
		path = filepath.Join(basepath, "iium_2024_2025_1.json")
	}

	catalog, err := os.ReadFile(path)
	if err != nil {
		logger.Error(err)
	}

	_catalog := json.RawMessage(catalog)

	return _catalog, nil
}
