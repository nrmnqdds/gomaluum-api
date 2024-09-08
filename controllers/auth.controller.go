package controllers

import (
	"net/http"

	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/internal"
	"github.com/nrmnqdds/gomaluum-api/services/auth"

	"github.com/labstack/echo/v4"
)

// @Title LoginHandler
// @Description Login to i-Ma'luum
// @Tags login
// @Accept json
// @Produce json
// @Param user body dtos.LoginDTO true "User object"
// @Success 200 {object} map[string]interface{}
// @Router /api/login [post]
func LoginHandler(c echo.Context) error {
	user := dtos.LoginDTO{}

	logger := internal.NewLogger()

	if c.Bind(&user) != nil {

		logger.Error("Invalid request payload!")

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
		logger.Error("Invalid request payload!", validationErr)
		return c.JSON(http.StatusBadRequest, response)
	}

	data, err := auth.LoginUser(&user)
	if err != nil {
		response := dtos.Response{
			Status:  err.StatusCode,
			Message: err.Message,
			Data:    nil,
		}
		logger.Error("Invalid request payload!", err)
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
