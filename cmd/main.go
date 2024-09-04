package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nrmnqdds/gomaluum-api/handlers/auth"
	"github.com/nrmnqdds/gomaluum-api/handlers/scraper"
	slogecho "github.com/samber/slog-echo"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(slogecho.New(logger))
	e.Use(middleware.Recover())

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.POST("/api/v1/login", auth.LoginUser)

	e.GET("/api/v1/profile", scraper.ProfileScraper)

	e.GET("/api/v1/schedule", scraper.ScheduleScraper)

	e.Logger.Fatal(e.Start(":1323"))
}
