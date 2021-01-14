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
}

func NewHealthChecker(o HealthCheckerOptions) hexa.HealthChecker {
	return &echoHealthChecker{
		echo:              o.Echo,
		livenessRoute:     o.LivenessRoute,
		readinessRoute:    o.ReadinessRoute,
		statusRoute:       o.StatusRoute,
		statusMiddlewares: o.StatusMiddlewares,
	}
}

func (c *echoHealthChecker) StartServer(r hexa.HealthReporter) error {
	if c.initialized {
		return tracer.Trace(errors.New("you can not start server after first start, create new instance"))
	}

	if c.livenessRoute != "" {
		c.echo.GET(c.livenessRoute, c.checkLiveness(r))
	}

	if c.readinessRoute != "" {
		c.echo.GET(c.readinessRoute, c.checkReadiness(r))
	}

	if c.statusRoute != "" {
		c.echo.GET(c.statusRoute, c.checkStatus(r), c.statusMiddlewares...)
	}

	// This healthChecker don't start the Echo server. If you need to start the Echo
	// server by healthChecker itself, implement another please.
	return nil
}

func (c *echoHealthChecker) StopServer() error {
	// Don't need to do anything.
	return nil
}

func (c *echoHealthChecker) checkLiveness(r hexa.HealthReporter) echo.HandlerFunc {
	return func(c echo.Context) error {
		status := r.LivenessStatus(context.Background())
		c.Response().Header().Set(hexa.LivenessStatusKey, string(status))

		if status != hexa.StatusAlive {
			return livenessReply.SetHTTPStatus(http.StatusInternalServerError)
		}

		return livenessReply
	}
}
func (c *echoHealthChecker) checkReadiness(r hexa.HealthReporter) echo.HandlerFunc {
	return func(c echo.Context) error {
		status := r.ReadinessStatus(context.Background())
		c.Response().Header().Set(hexa.ReadinessStatusKey, string(status))

		if status != hexa.StatusReady {
			return readinessReply.SetHTTPStatus(http.StatusServiceUnavailable)
		}

		return readinessReply
	}
}

func (c *echoHealthChecker) checkStatus(r hexa.HealthReporter) echo.HandlerFunc {
	return func(c echo.Context) error {
		report := r.HealthReport(context.Background())

		c.Response().Header().Set(hexa.LivenessStatusKey, string(report.Alive))
		c.Response().Header().Set(hexa.ReadinessStatusKey, string(report.Ready))

		return StatusReply.SetData(gutil.StructToMap(report))
	}
}

func DefaultHealthCheckerOptions(echo *echo.Echo, statusMiddlewares ...echo.MiddlewareFunc) HealthCheckerOptions {
	return HealthCheckerOptions{
		Echo:              echo,
		LivenessRoute:     "/alive",
		ReadinessRoute:    "/ready",
		StatusRoute:       "/status",
		StatusMiddlewares: statusMiddlewares,
	}
}
