package controllers

import (
	"net/http"

	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/internal"
	"github.com/nrmnqdds/gomaluum-api/services/scraper"

	"github.com/labstack/echo/v4"
)

// @Title GetResultHandler
// @Description Get result from i-Ma'luum
// @Tags scraper
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/result [get]
func GetResultHandler(c echo.Context) error {
	data, err := scraper.ResultScraper(c)

	logger := internal.NewLogger()

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
		Message: "Successfully get user result!",
		Data:    data,
	}

	return c.JSON(http.StatusOK, response)
}
