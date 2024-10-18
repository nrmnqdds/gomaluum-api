package controllers

import (
	"net/http"

	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/helpers"
	"github.com/nrmnqdds/gomaluum-api/services/auth"
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
	logger, _ := helpers.NewLogger()

	result := dtos.ScheduleRequestProps{
		Echo: c,
	}

	data, err := scraper.ResultScraper(&result)
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

func PostResultHandler(c echo.Context) error {
	user := dtos.LoginDTO{}
	logger, _ := helpers.NewLogger()

	if c.Bind(&user) != nil {

		logger.Error("Invalid request payload!")

		response := dtos.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request payload!",
			Data:    nil,
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	if validationErr := helpers.Validator.Struct(&user); validationErr != nil {
		response := dtos.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request payload!",
			Data:    nil,
		}
		logger.Error("Invalid request payload!", validationErr)
		return c.JSON(http.StatusBadRequest, response)
	}

	loginRes, err := auth.LoginUser(&user)
	if err != nil {
		response := dtos.Response{
			Status:  err.StatusCode,
			Message: err.Message,
			Data:    nil,
		}
		logger.Error(err)
		return c.JSON(http.StatusInternalServerError, response)
	}

	result := dtos.ScheduleRequestProps{
		Echo:  c,
		Token: loginRes.Token,
	}

	data, err := scraper.ResultScraper(&result)
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
