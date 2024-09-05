package dtos

import "github.com/labstack/echo/v4"

type Response struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    *echo.Map `json:"data"`
}

