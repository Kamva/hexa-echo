package kecho

import (
	"github.com/Kamva/kitty"
	"github.com/labstack/echo/v4"
)

type (
	// UserFinder is a function to use to find current user by jwt claims.
	UserFinder func(jwtClaim interface{}) kitty.User

	// CurrentUserConfig is the config to use in CurrentUser middleware.
	CurrentUserConfig struct {
		UserFinder     UserFinder
		UserContextKey string
		JWTContextKey  string
	}
)

var (
	// CurrentUserContextKey is the context key to set
	// the current user in the request context.
	CurrentUserContextKey = "user"
)

func CurrentUser(userFinder UserFinder) echo.MiddlewareFunc {
	return CurrentUserWithConfig(CurrentUserConfig{
		UserFinder:     userFinder,
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
			ctx.Set(cfg.UserContextKey, cfg.UserFinder(jwt))

			return next(ctx)
		}
	}
}
