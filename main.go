package main

import (
	"embed"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nrmnqdds/gomaluum-api/controllers"
	services "github.com/nrmnqdds/gomaluum-api/services/catalogs"
)

//go:embed docs/*
var DocsPath embed.FS

func main() {
	f, err := os.Create("cpu.pprof")
	if err != nil {
		panic(err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	if os.Getenv("APP_ENV") == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error reading .env file!")
		}
		log.Println("Loaded .env file")
	}
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

	// e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))

	// This middleware is used to recover from panics anywhere in the chain, log the panic (and a stack trace), and set a status code of 500.
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/reference")
	})

	e.GET("/reference", func(c echo.Context) error {
		swaggerPath, err := DocsPath.ReadFile("docs/swagger/swagger.json")
		if err != nil {
			log.Fatalf("could not read swagger.json: %v", err)
		}
		scalarMetadataPath, err := DocsPath.ReadFile("docs/scalar-metadata.json")
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
	services.DocsPath = DocsPath
	g.GET("/catalog", controllers.CatalogHandler)

	// Ads
	g.GET("/ads", controllers.AdsHandler)

	PORT := "1323"

	if os.Getenv("PORT") != "" {
		PORT = os.Getenv("PORT")
	}

	if err := e.Start(":" + PORT); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
	log.Println("App Server started")
}
