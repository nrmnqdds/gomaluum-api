package controllers

import (
	"net/http"

	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/internal"
	"github.com/nrmnqdds/gomaluum-api/services/scraper"

	"github.com/labstack/echo/v4"
)

func GetProfileHandler(c echo.Context) error {
	user := dtos.LoginDTO{}

	if c.Bind(&user) != nil {
		response := dtos.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request payload!",
			Data:    nil,
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	if validationErr := internal.Validator.Struct(&user); validationErr != nil {
		response := dtos.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request payload!",
			Data:    nil,
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	data, err := scraper.ProfileScraper(c)
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
		Message: "Successfully get user profile!",
		Data:    &echo.Map{"data": data},
	}

	return c.JSON(http.StatusCreated, response)
}
