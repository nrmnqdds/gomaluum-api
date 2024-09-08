package controllers

import (
	"net/http"

	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/services/scraper"

	"github.com/labstack/echo/v4"
)

// @Title GetScheduleHandler
// @Description Get schedule from i-Ma'luum
// @Tags scraper
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/schedule [post]
func GetScheduleHandler(c echo.Context) error {
	data, err := scraper.ScheduleScraper(c)
	if err != nil {
		response := dtos.Response{
			Status:  err.StatusCode,
			Message: err.Message,
			Data:    nil,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := dtos.Response{
		Status:  http.StatusOK,
		Message: "Successfully get user schedule!",
		Data:    &echo.Map{"schedule": data},
	}

	return c.JSON(http.StatusOK, response)
}
