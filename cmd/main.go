package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"

	otelmid "go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"golang.org/x/time/rate"

	"google.golang.org/grpc/credentials"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/maruel/panicparse/v2/stack"
	"github.com/nrmnqdds/gomaluum-api/controllers"
	_ "github.com/nrmnqdds/gomaluum-api/docs/swagger"
	"github.com/nrmnqdds/gomaluum-api/helpers"
	echoSwagger "github.com/swaggo/echo-swagger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	serviceName  = helpers.GetEnv("SERVICE_NAME")
	collectorURL = helpers.GetEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = helpers.GetEnv("INSECURE_MODE")
)

func initTracer() func(context.Context) error {
	var secureOption otlptracegrpc.Option
	logger, _ := helpers.NewLogger()

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
		logger.Errorf("Failed to create exporter: %v", err)
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
	logger, _ := helpers.NewLogger()

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

	cleanup := initTracer()
	defer func() {
		if err := cleanup(context.Background()); err != nil {
			logger.Errorf("Failed to shutdown exporter: %v", err)
		}
	}()

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

	// Set up OpenTelemetry middleware
	e.Use(otelmid.Middleware(serviceName))

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

	g.GET("/catalog", controllers.CatalogHandler)

	g.GET("/ads", controllers.AdsHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
