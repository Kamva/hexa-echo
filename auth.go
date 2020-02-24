package kecho

import (
	"github.com/Kamva/kitty"
	"github.com/labstack/echo/v4"
)

type GateMiddlewareConfig struct {
	MustBeGuest bool
}

// GuestMiddleware is a middleware that force user to be guest to access to specific API.
// GuestMiddleware should be after the kittyContext middleware.
func GuestMiddleware() echo.MiddlewareFunc {
	return UserGateMiddleware(GateMiddlewareConfig{MustBeGuest: true})
}

// AuthMiddleware is a middleware that force user to authenticate to access to specific API.
// AuthMiddleware should be after the kittyContext middleware.
func AuthMiddleware() echo.MiddlewareFunc {
	return UserGateMiddleware(GateMiddlewareConfig{MustBeGuest: false})
}

// UserGateMiddleware is a middleware to specify user should be authenticated or
// be guest to access to specific API.
func UserGateMiddleware(cfg GateMiddlewareConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			kittyCtx := ctx.Get(ContextKeyKittyCtx).(kitty.Context)
			u := kittyCtx.User()

			// validate guest rule:
			if cfg.MustBeGuest && !u.IsGuest() {
				return errUserMustBeGuest
			}

			if u.IsGuest() {
				return errUserNeedToAuthenticate
			}

			return next(ctx)
		}
	}
}
