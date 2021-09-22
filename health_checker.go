package hecho

import (
	"context"
	"errors"
	"net/http"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
)

var livenessReply = hexa.NewReply(http.StatusOK, "liveness")
var readinessReply = hexa.NewReply(http.StatusOK, "readiness")
var StatusReply = hexa.NewReply(http.StatusOK, "status")

type HealthCheckerOptions struct {
	Echo           *echo.Echo
	LivenessRoute  string // Empty value means disabled liveness
	ReadinessRoute string // Empty route means disabled readiness
	StatusRoute    string // Empty route means disabled status

	StatusMiddlewares []echo.MiddlewareFunc
	Reporter          hexa.HealthReporter
	Logger            hexa.Logger
	Translator        hexa.Translator
}

// this healthChecker assume echo server is your default server and start,stop
// of it is with yourself.
// Notes:
// - This healthChecker don't start the Echo server. it should start by yourself.
// - You must set hecho HTTPErrorHandler as echo error handler. that
//   is because we use replies in this healthChecker and need to convert
//   replies to response.
type echoHealthChecker struct {
	initialized    bool
	echo           *echo.Echo
	livenessRoute  string
	readinessRoute string
	statusRoute    string

	statusMiddlewares []echo.MiddlewareFunc
	reporter          hexa.HealthReporter

	l hexa.Logger
	t hexa.Translator
}

func NewHealthChecker(o HealthCheckerOptions) hexa.HealthChecker {
	return &echoHealthChecker{
		echo:              o.Echo,
		livenessRoute:     o.LivenessRoute,
		readinessRoute:    o.ReadinessRoute,
		statusRoute:       o.StatusRoute,
		statusMiddlewares: o.StatusMiddlewares,
		reporter:          o.Reporter,
	}
}

func (h *echoHealthChecker) Run() error {
	if h.initialized {
		return tracer.Trace(errors.New("you can not start server after first start, create new instance"))
	}

	if h.livenessRoute != "" {
		h.echo.GET(h.livenessRoute, h.checkLiveness())
	}

	if h.readinessRoute != "" {
		h.echo.GET(h.readinessRoute, h.checkReadiness())
	}

	if h.statusRoute != "" {
		h.echo.GET(h.statusRoute, h.checkStatus(), h.statusMiddlewares...)
	}

	// This healthChecker don't start the Echo server. If you need to start the Echo
	// server by healthChecker itself, implement another please.
	return nil
}

func (h *echoHealthChecker) Shutdown(_ context.Context) error {
	// Don't need to do anything.
	return nil
}

func (h *echoHealthChecker) checkLiveness() echo.HandlerFunc {
	return func(c echo.Context) error {
		status := h.reporter.LivenessStatus(context.Background())
		c.Response().Header().Set(hexa.LivenessStatusKey, string(status))

		if status != hexa.StatusAlive {
			return WriteWithOpts(c, h.l, h.t, livenessReply.SetHTTPStatus(http.StatusInternalServerError))
		}

		return WriteWithOpts(c, h.l, h.t, livenessReply)
	}
}
func (h *echoHealthChecker) checkReadiness() echo.HandlerFunc {
	return func(c echo.Context) error {
		status := h.reporter.ReadinessStatus(context.Background())
		c.Response().Header().Set(hexa.ReadinessStatusKey, string(status))

		if status != hexa.StatusReady {
			return WriteWithOpts(c, h.l, h.t, readinessReply.SetHTTPStatus(http.StatusServiceUnavailable))
		}

		return WriteWithOpts(c, h.l, h.t, readinessReply)
	}
}

func (h *echoHealthChecker) checkStatus() echo.HandlerFunc {
	return func(c echo.Context) error {
		report := h.reporter.HealthReport(context.Background())

		c.Response().Header().Set(hexa.LivenessStatusKey, string(report.Alive))
		c.Response().Header().Set(hexa.ReadinessStatusKey, string(report.Ready))
		return WriteWithOpts(c, h.l, h.t, StatusReply.SetData(gutil.StructToMap(report)))
	}
}

func DefaultHealthCheckerOptions(echo *echo.Echo, r hexa.HealthReporter, l hexa.Logger, t hexa.Translator, statusMiddlewares ...echo.MiddlewareFunc) HealthCheckerOptions {
	return HealthCheckerOptions{
		Echo:              echo,
		LivenessRoute:     "/live",
		ReadinessRoute:    "/ready",
		StatusRoute:       "/status",
		StatusMiddlewares: statusMiddlewares,
		Reporter:          r,
		Logger:            l,
		Translator:        t,
	}
}
