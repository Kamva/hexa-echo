package main

import (
	"context"
	"net/http"

	"github.com/kamva/gutil"
	hecho "github.com/kamva/hexa-echo"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/hexa/htel"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	service     = "hexa-demo"
	environment = "dev"
	id          = 1
)

var ot htel.OpenTelemetry

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func tracerProvider(url string) (*tracesdk.TracerProvider, error) {

	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tp, nil
}

func main() {
	// Run jaeger server:
	// docker run -d --name jaeger \
	//  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
	//  -p 5775:5775/udp \
	//  -p 6831:6831/udp \
	//  -p 6832:6832/udp \
	//  -p 5778:5778 \
	//  -p 16686:16686 \
	//  -p 14268:14268 \
	//  -p 14250:14250 \
	//  -p 9411:9411 \
	//  jaegertracing/all-in-one:1.26
	// Navigate to http://localhost:16686/ to view jaeger UI.
	tp, err := tracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		gutil.PanicErr(err)
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	defer func() {
		tp.Shutdown(context.Background())
	}()

	l := hlog.NewPrinterDriver(hlog.DebugLevel)
	e := echo.New()
	e.Debug = true

	e.Logger = hecho.HexaToEchoLogger(l, "debug")

	e.Use(hecho.Tracing(hecho.TracingConfig{
		Propagator:     propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
		TracerProvider: tp,
		ServerName:     "lab",
	}))

	e.Use(hecho.Recover())

	e.HTTPErrorHandler = hecho.HTTPErrorHandler(l, hexatranslator.NewEmptyDriver(), true)
	e.GET("/hi", func(c echo.Context) error {
		//var a map[string]interface{}
		//a["a"] = "12"
		c.String(http.StatusAccepted, "hi :)")
		//return errors.New("fake error")
		return nil
	})

	e.Logger.Fatal(e.Start(":4444"))
}
