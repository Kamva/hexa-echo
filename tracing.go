package hecho

import (
	"strings"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/htel"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type TracingConfig struct {
	Propagator propagation.TextMapPropagator
	Tracer     trace.Tracer
	ServerName string
	SpanName   string
}

// Tracing Enables distributed tracing using openTelemetry library.
// In echo if a handler panic error, it will catch by the `Recover`
// middleware, to get panic errors too, please use this middleware
// before the Recover middleware, so it will get the recovered
// errors too.
// You can use TracingDataFromUserContext middleware to set user_id
// and correlation_id too.
func Tracing(cfg TracingConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			r := context.Request()
			attrs := semconv.HTTPServerAttributesFromHTTPRequest(cfg.ServerName, r.URL.Path, r)

			// Extract the parent from the request, but this is a gateway that users
			// send request to it, check if propagation from external requests has any
			// security issue.
			// TODO: check if we should not get parent for external requests, remove it.
			ctx := cfg.Propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			ctx, span := cfg.Tracer.Start(ctx, cfg.SpanName, trace.WithAttributes(attrs...))
			context.SetRequest(r.Clone(ctx))

			defer func() {
				span.End()
			}()

			err := next(context)
			if isInternalErr(err) { // ignore hexa.Reply or hexa.Error with code out of 5XX range.
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}

			// Set http status code and user's data as context:
			semconv.HTTPAttributesFromHTTPStatusCode(context.Response().Status)

			return err
		}
	}
}

// TracingDataFromUserContext sets some tags,... on tracing span
// using hexa context. This middleware should be after hexa context
// middleware because if needs to the hexa context.
func TracingDataFromUserContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			hexaCtx, ok := c.Get(ContextKeyHexaCtx).(hexa.Context)
			if !ok {
				return next(c)
			}

			user := hexaCtx.User()
			// Add user's id, correlation_id
			span := trace.SpanFromContext(hexaCtx)
			span.SetAttributes(
				semconv.EnduserIDKey.String(user.Identifier()), // enduser.id
				htel.EnduserUsernameKey.String(user.Username()), // enduser.username
				semconv.EnduserRoleKey.String(strings.Join(user.Roles(), ",")), // enduser.role
				htel.CorrelationIDKey.String(hexaCtx.CorrelationID()), // ctx.correlation_id
			)

			return next(c)
		}
	}
}
