package controllers

import (
	"net/http"

	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/helpers"
	"github.com/nrmnqdds/gomaluum-api/services/scraper"

	"github.com/labstack/echo/v4"
)

// @Title AdsHandler
// @Description SOUQ Ads
// @Tags ads
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/ads [post]
func AdsHandler(c echo.Context) error {
	logger, _ := helpers.NewLogger()

	data, err := scraper.AdsScraper()
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
		Message: "Successfully get ads!",
		Data:    data,
	}

	return c.JSON(http.StatusOK, response)
}
