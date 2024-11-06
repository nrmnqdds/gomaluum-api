package main

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nrmnqdds/gomaluum-api/controllers"
	"github.com/nrmnqdds/gomaluum-api/helpers"
)

// @title Gomaluum API
// @version 1.0
// @description This is a simple API for Gomaluum project.
func main() {
	e := echo.New()

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
		Output: e.Logger.Output(),
	}))

	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))

	// This middleware is used to recover from panics anywhere in the chain, log the panic (and a stack trace), and set a status code of 500.
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/reference")
	})

	e.GET("/reference", func(c echo.Context) error {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			// SpecURL: "https://generator3.swagger.io/openapi.json",// allow external URL or local path file
			SpecURL: helpers.OpenAPISpecPath,
			CustomOptions: scalar.CustomOptions{
				PageTitle: "GoMaluum API Reference",
			},
			DarkMode: true,
		})
		if err != nil {
			fmt.Printf("%v", err)
		}

		htmlBlob := []byte(htmlContent)

		return c.HTMLBlob(http.StatusOK, htmlBlob)
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	g := e.Group("/api")

	g.POST("/login", controllers.LoginHandler)

	g.GET("/profile", controllers.GetProfileHandler)

	// Schedule
	g.GET("/schedule", controllers.GetScheduleHandler)
	g.POST("/schedule", controllers.PostScheduleHandler)

	g.GET("/result", controllers.GetResultHandler)
	g.POST("/result", controllers.PostResultHandler)

	// Catalog
	g.GET("/catalog", controllers.CatalogHandler)

	// Ads
	g.GET("/ads", controllers.AdsHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
