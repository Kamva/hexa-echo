package main

import (
	"fmt"
	"log"
	"net/http"

	hecho "github.com/kamva/hexa-echo"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/hexa/htel"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	service     = "hexa-demo"
	environment = "dev"
	id          = 1
)

var ot htel.OpenTelemetry

func initMeter() {
	config := prometheus.Config{}
	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
		controller.WithResource(resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id))),
	)

	exporter, err := prometheus.New(config, c)

	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}
	global.SetMeterProvider(exporter.MeterProvider())

	http.HandleFunc("/", exporter.ServeHTTP)
	go func() {
		_ = http.ListenAndServe(":2222", nil)
	}()

	fmt.Println("Prometheus server running on :2222")
}

func main() {

	initMeter()

	l := hlog.NewPrinterDriver(hlog.DebugLevel)
	e := echo.New()
	e.Debug = true

	e.Logger = hecho.HexaToEchoLogger(l, "debug")

	e.Use(hecho.Metrics(hecho.MetricsConfig{
		MeterProvider: global.GetMeterProvider(),
		ServerName:    "lab",
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
