package kecho

import (
	"errors"
	"github.com/Kamva/kitty"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type (
	// UserFinderByJwtSub is a function to use to find current user by jwt claims.
	UserFinderByJwtSub func(sub string) (kitty.User, error)

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

// CurrentUser is a middleware to set the user in the context.
// If provided jwt, so this function find user and set it as user
// otherwise set guest user.
func CurrentUser(userFinder UserFinderByJwtSub) echo.MiddlewareFunc {
	return CurrentUserWithConfig(CurrentUserConfig{
		UserFinderByJwtSub: userFinder,
		UserContextKey:     CurrentUserContextKey,
		JWTContextKey:      JwtContextKey,
	})
}

// CurrentUser is a middleware to set the user in the context.
// If provided jwt, so this function find user and set it as user
// otherwise set guest user.
func CurrentUserWithConfig(cfg CurrentUserConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// Get jwt (if exists)
			if token, ok := ctx.Get(cfg.JWTContextKey).(*jwt.Token); ok {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					// Set the user.
					user, err := cfg.UserFinderByJwtSub(claims["sub"].(string))

					if err != nil {
						return err
					}

					ctx.Set(cfg.UserContextKey, user)

					return next(ctx)
				}

				return errors.New("JWT claims is not valid")
			}

			ctx.Set(cfg.UserContextKey, kitty.NewGuestUser())
			return next(ctx)
		}
	}
}
