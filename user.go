package kecho

import (
	"github.com/Kamva/kitty"
	"github.com/labstack/echo/v4"
)

type (
	// UserGetter is a function to use to get current user from jwt claims.
	UserGetter func(jwtClaim interface{}) kitty.User

	// CurrentUserConfig is the config to use in CurrentUser middleware.
	CurrentUserConfig struct {
		UserGetter     UserGetter
		UserContextKey string
		JWTContextKey  string
	}
)

var (
	// CurrentUserContextKey is the context key to set
	// the current user in the request context.
	CurrentUserContextKey = "user"
)

func CurrentUser(getter UserGetter) echo.MiddlewareFunc {
	return CurrentUserWithConfig(CurrentUserConfig{
		UserGetter:     getter,
		UserContextKey: CurrentUserContextKey,
		JWTContextKey:  JwtContextKey,
	})
}

// CurrentUser is a middleware to set the user in the context.
// be guest to access to specific API.
func CurrentUserWithConfig(cfg CurrentUserConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// Get jwt
			jwt := ctx.Get(cfg.JWTContextKey)

			// Set the user.
			ctx.Set(cfg.UserContextKey, cfg.UserGetter(jwt))

			return next(ctx)
		}
	}
}
