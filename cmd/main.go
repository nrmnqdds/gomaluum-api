package main

import (
	"net/http"

	"github.com/labstack/echo-contrib/echoprometheus"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nrmnqdds/gomaluum-api/controllers"
	_ "github.com/nrmnqdds/gomaluum-api/docs/swagger"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Gomaluum API
// @version 1.0
// @description This is a simple API for Gomaluum project.
func main() {
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
		Output: e.Logger.Output(),
	}))

	// This middleware is used to recover from panics anywhere in the chain, log the panic (and a stack trace), and set a status code of 500.
	e.Use(middleware.Recover())
	e.Use(echoprometheus.NewMiddleware("gomaluum")) // adds middleware to gather metrics

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics

	e.POST("/api/login", controllers.LoginHandler)

	e.GET("/api/profile", controllers.GetProfileHandler)

	e.GET("/api/schedule", controllers.GetScheduleHandler)

	e.GET("/api/result", controllers.GetResultHandler)

	e.GET("/api/catalog", controllers.CatalogHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
