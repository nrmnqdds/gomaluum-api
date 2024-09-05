package controllers

import (
	"net/http"

	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/internal"
	"github.com/nrmnqdds/gomaluum-api/services/auth"

	"github.com/labstack/echo/v4"
)

func LoginHandler(c echo.Context) error {
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

	data, err := auth.LoginUser(&user)
	if err != nil {
		response := dtos.Response{
			Status:  err.StatusCode,
			Message: err.Message,
			Data:    nil,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	c.SetCookie(&http.Cookie{
		Name:  "MOD_AUTH_CAS",
		Value: data.Token,
	})

	response := dtos.Response{
		Status:  http.StatusOK,
		Message: "Successfully login!",
		Data:    &echo.Map{"data": data},
	}

	return c.JSON(http.StatusCreated, response)
}
