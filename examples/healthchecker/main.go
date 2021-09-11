package main

import (
	"context"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	hecho "github.com/kamva/hexa-echo"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"github.com/labstack/echo/v4"
)

type HealthExample struct {
}

func (h *HealthExample) HealthIdentifier() string {
	return "health_example"
}

func (h *HealthExample) LivenessStatus(ctx context.Context) hexa.LivenessStatus {
	return hexa.StatusAlive
}

func (h *HealthExample) ReadinessStatus(ctx context.Context) hexa.ReadinessStatus {
	return hexa.StatusReady
}

func (h *HealthExample) HealthStatus(ctx context.Context) hexa.HealthStatus {
	return hexa.HealthStatus{
		Id:    h.HealthIdentifier(),
		Tags:  map[string]string{"I'm": "ok :)"},
		Alive: h.LivenessStatus(ctx),
		Ready: h.ReadinessStatus(ctx),
	}
}

var l = hlog.NewPrinterDriver(hlog.DebugLevel)

func main() {
	e := echo.New()
	e.Debug = true
	e.Logger = hecho.HexaToEchoLogger(l, "debug")

	// hecho ErrorHandler is required for healthChecker(to convert reply to response)
	e.HTTPErrorHandler = hecho.HTTPErrorHandler(l, hexatranslator.NewEmptyDriver(), true)

	r := hexa.NewHealthReporter()
	r.AddToChecks(&HealthExample{})
	gutil.PanicErr(hecho.NewHealthChecker(hecho.DefaultHealthCheckerOptions(e, r)).Run())
	e.Logger.Fatal(e.Start(":4444"))
}
