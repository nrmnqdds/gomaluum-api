package services

import (
	"embed"
	"encoding/json"

	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/helpers"
)

var DocsPath embed.FS

func CatalogScraper() (interface{}, *dtos.CustomError) {
	logger, _ := helpers.NewLogger()

	catalog, err := DocsPath.ReadFile("docs/iium_2024_2025_1.json")
	if err != nil {
		logger.Error(err)
	}

	_catalog := json.RawMessage(catalog)

	return _catalog, nil
}
