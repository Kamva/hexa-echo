package hecho

import (
	"github.com/kamva/hexa"
	"github.com/labstack/echo/v4"
)

// SetContextLogger set the hexa logger on each context.
func SetContextLogger(level string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// Set context logger
			hexaContext := ctx.Get(ContextKeyHexaCtx).(hexa.Context)
			ctx.SetLogger(HexaToEchoLogger(hexaContext.Logger(), level))
			return next(ctx)
		}
	}
}
