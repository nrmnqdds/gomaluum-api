package controllers

import (
	"net/http"

	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/helpers"
	"github.com/nrmnqdds/gomaluum-api/services/catalogs"

	"github.com/labstack/echo/v4"
)

// @Title CatalogHandler
// @Description Get catalog
// @Tags catalog
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/catalog [get]
func CatalogHandler(c echo.Context) error {
	logger, _ := helpers.NewLogger()

	data, err := catalog.CatalogScraper(c)
	if err != nil {
		response := dtos.Response{
			Status:  err.StatusCode,
			Message: err.Message,
			Data:    nil,
		}
		logger.Error(err)
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := dtos.Response{
		Status:  http.StatusOK,
		Message: "Successfully get catalog!",
		Data:    data,
	}

	return c.JSON(http.StatusCreated, response)
}
