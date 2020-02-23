package kecho

import (
	"errors"
	"github.com/Kamva/kitty"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type (
	// UserFinderByJwtSub is a function to use to find current user by jwt claims.
	UserFinderByJwtSub func(sub string) kitty.User

	// CurrentUserConfig is the config to use in CurrentUser middleware.
	CurrentUserConfig struct {
		UserFinderByJwtSub UserFinderByJwtSub
		UserContextKey     string
		JWTContextKey      string
	}
)

var (
	// CurrentUserContextKey is the context key to set
	// the current user in the request context.
	CurrentUserContextKey = "user"
)

func CurrentUser(userFinder UserFinderByJwtSub) echo.MiddlewareFunc {
	return CurrentUserWithConfig(CurrentUserConfig{
		UserFinderByJwtSub: userFinder,
		UserContextKey:     CurrentUserContextKey,
		JWTContextKey:      JwtContextKey,
	})
}

// CurrentUser is a middleware to set the user in the context.
// be guest to access to specific API.
func CurrentUserWithConfig(cfg CurrentUserConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// Get jwt
			token := ctx.Get(cfg.JWTContextKey).(*jwt.Token)

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// Set the user.
				ctx.Set(cfg.UserContextKey, cfg.UserFinderByJwtSub(claims["sub"].(string)))

				return next(ctx)
			}

			return errors.New("JWT claims is not valid")
		}
	}
}
