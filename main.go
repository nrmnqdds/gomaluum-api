package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nrmnqdds/gomaluum-api/controllers"
	"golang.org/x/time/rate"
)

//go:embed docs/*
var SwaggerDocsPath embed.FS

var echoLambda *echoadapter.EchoLambda

func main() {
	lambda.Start(Handler)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error reading .env file!")
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

	isLambda := os.Getenv("LAMBDA")

	if isLambda == "TRUE" {
		// echoadapter := echoadapter.New(e)
		//
		// lambda.Start(echoadapter.Proxy)
		// lambdaAdapter := &LambdaAdapter{Echo: e}
		// lambda.Start(lambdaAdapter.Handler)

		echoLambda = echoadapter.New(e)
	} else {
		PORT := "1323"

		if os.Getenv("PORT") != "" {
			PORT = os.Getenv("PORT")
		}

		if err := e.Start(":" + PORT); err != nil {
			log.Fatalf("could not start server: %v", err)
		}
		log.Println("App Server started")
	}
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return echoLambda.ProxyWithContext(ctx, req)
}
