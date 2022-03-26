package hecho

import (
	"github.com/kamva/hexa/hlog"
	"github.com/labstack/echo/v4"
)

// SetContextLogger set the hexa logger on each context.
func SetContextLogger(level string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// Set context logger
			ctx.SetLogger(HexaToEchoLogger(hlog.CtxLogger(ctx.Request().Context()), level))
			return next(ctx)
		}
	}
}
