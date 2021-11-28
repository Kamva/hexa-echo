package hecho

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type MetricsConfig struct {
	Skipper       middleware.Skipper
	MeterProvider metric.MeterProvider
	ServerName    string
}

func Metrics(cfg MetricsConfig) echo.MiddlewareFunc {
	if cfg.Skipper == nil {
		cfg.Skipper = middleware.DefaultSkipper
	}

	meter := metric.Must(cfg.MeterProvider.Meter(instrumentationName))
	requestCounter := meter.NewFloat64Counter("requests_total")
	requestDuration := meter.NewFloat64Histogram("requests_duration_second")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.Skipper(c) {
				return next(c)
			}

			startTime := time.Now()
			r := c.Request()

			// Extract the parent from the request, but this is a gateway that users
			// send request to it, check if propagation from external requests has any
			// security issue.
			attrs := semconv.NetAttributesFromHTTPRequest("tcp", r)
			attrs = append(attrs, semconv.EndUserAttributesFromHTTPRequest(r)...)
			attrs = append(attrs, semconv.HTTPServerAttributesFromHTTPRequest(cfg.ServerName, c.Path(), r)...)

			err := next(c)
			if err != nil {
				c.Error(err) // apply the error to set the response code
			}

			attrs = append(attrs, semconv.HTTPAttributesFromHTTPStatusCode(c.Response().Status)...)

			elapsed := float64(time.Since(startTime)) / float64(time.Second)

			requestCounter.Add(r.Context(), 1, attrs...)
			requestDuration.Record(r.Context(), elapsed, attrs...)

			return nil // we applied the error, so we don't need to return it again.
		}
	}
}
