package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"

	"golang.org/x/time/rate"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/maruel/panicparse/v2/stack"
	"github.com/nrmnqdds/gomaluum-api/controllers"
	_ "github.com/nrmnqdds/gomaluum-api/docs/swagger"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Gomaluum API
// @version 1.0
// @description This is a simple API for Gomaluum project.
func main() {
	e := echo.New()

	parseStack := func(rawStack []byte) stack.Stack {
		s, _, err := stack.ScanSnapshot(bytes.NewReader(rawStack), io.Discard, stack.DefaultOpts())
		if err != nil && err != io.EOF {
			panic(err)
		}

		if len(s.Goroutines) > 1 {
			panic(errors.New("provided stacktrace had more than one goroutine"))
		}

		return s.Goroutines[0].Signature.Stack
	}

	parsedStack := parseStack(debug.Stack())
	fmt.Printf("parsedStack: %#v", parsedStack)

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
		return c.String(http.StatusOK, "You look lost. Try /swagger/index.html")
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

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
