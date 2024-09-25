package main

import (
	"context"
	"net/http"
	"strings"

	otelmid "go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

	"google.golang.org/grpc/credentials"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nrmnqdds/gomaluum-api/controllers"
	_ "github.com/nrmnqdds/gomaluum-api/docs/swagger"
	"github.com/nrmnqdds/gomaluum-api/internal"
	echoSwagger "github.com/swaggo/echo-swagger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	serviceName  = internal.GetEnv("SERVICE_NAME")
	collectorURL = internal.GetEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = internal.GetEnv("INSECURE_MODE")
	logger       = internal.NewLogger()
)

func initTracer() func(context.Context) error {
	var secureOption otlptracegrpc.Option

	if strings.ToLower(insecure) == "false" || insecure == "0" || strings.ToLower(insecure) == "f" {
		secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)
	if err != nil {
		logger.Fatalf("Failed to create exporter: %v", err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		logger.Fatalf("Could not set resources: %v", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}

// @title Gomaluum API
// @version 1.0
// @description This is a simple API for Gomaluum project.
func main() {
	e := echo.New()

	cleanup := initTracer()
	defer func() {
		if err := cleanup(context.Background()); err != nil {
			logger.Fatalf("Failed to shutdown exporter: %v", err)
		}
	}()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
		Output: e.Logger.Output(),
	}))

	// Set up OpenTelemetry middleware
	e.Use(otelmid.Middleware(serviceName))

	// This middleware is used to recover from panics anywhere in the chain, log the panic (and a stack trace), and set a status code of 500.
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "You look lost. Try /swagger/index.html")
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.POST("/api/login", controllers.LoginHandler)

	e.GET("/api/profile", controllers.GetProfileHandler)

	e.GET("/api/schedule", controllers.GetScheduleHandler)

	e.GET("/api/result", controllers.GetResultHandler)

	e.GET("/api/catalog", controllers.CatalogHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
