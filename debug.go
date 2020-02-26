package kecho

import (
	"github.com/labstack/echo/v4"
)

// KittyContext set kitty context on each request.
func DebugMode(e *echo.Echo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if !e.Debug {
				return errRouteAvaialbeInDebugMode
			}

			return next(ctx)
		}
	}
}
