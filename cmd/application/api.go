package application

import (
	"embed"
	"log"
	"net/http"

	"golang.org/x/time/rate"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nrmnqdds/gomaluum-api/controllers"
	"github.com/nrmnqdds/gomaluum-api/helpers"
)

var SwaggerDocsPath *embed.FS

// @title Gomaluum API
// @version 1.0
// @description This is a simple API for Gomaluum project.
func StartEchoServer() {
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
		swaggerPath, err := SwaggerDocsPath.ReadFile("docs/swagger/swagger.json")
		if err != nil {
			log.Fatalf("could not read swagger.json: %v", err)
		}
		scalarMetadataPath, err := SwaggerDocsPath.ReadFile("docs/scalar-metadata.json")
		if err != nil {
			log.Fatalf("could not read scalar-metadata.json: %v", err)
		}

		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecContent: string(swaggerPath),
			CustomOptions: scalar.CustomOptions{
				PageTitle: "GoMaluum API Reference",
			},
			MetaData: string(scalarMetadataPath),
			DarkMode: true,
		})
		if err != nil {
			log.Fatalf("Error initializing scalar: %v", err)
		}

		htmlBlob := []byte(htmlContent)

		return c.HTMLBlob(http.StatusOK, htmlBlob)
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	mon := asynqmon.New(asynqmon.Options{
		RootPath: "/monitoring/tasks",
		RedisConnOpt: asynq.RedisClientOpt{
			Addr:     helpers.GetEnv("REDIS_URL"),
			Username: helpers.GetEnv("REDIS_USERNAME"),
			Password: helpers.GetEnv("REDIS_PASSWORD"),
		},
	})
	e.Any("/monitoring/tasks/*", echo.WrapHandler(mon))

	g := e.Group("/api")

	g.POST("/login", controllers.LoginHandler)

	g.GET("/profile", controllers.GetProfileHandler)

	// Schedule
	g.GET("/schedule", controllers.GetScheduleHandler)
	g.POST("/schedule", controllers.PostScheduleHandler)

	// Result
	g.GET("/result", controllers.GetResultHandler)
	g.POST("/result", controllers.PostResultHandler)

	// Catalog
	g.GET("/catalog", controllers.CatalogHandler)

	// Ads
	g.GET("/ads", controllers.AdsHandler)

	if err := e.Start(":" + helpers.GetEnv("PORT")); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
	log.Println("App Server started")
}
