package hecho

import (
	"errors"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type (
	// UserFinderBySub find the user by provided sub.
	UserFinderBySub func(sub string) (hexa.User, error)

	// CurrentUserConfig is the config to use in CurrentUser middleware.
	CurrentUserConfig struct {
		userSDK        hexa.UserSDK
		uf             UserFinderBySub // Can be nil if ExtendJWT is false.
		ExtendJWT      bool
		UserContextKey string
		JWTContextKey  string
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
func CurrentUser(uf UserFinderBySub, userSDK hexa.UserSDK) echo.MiddlewareFunc {
	return CurrentUserWithConfig(CurrentUserConfig{
		ExtendJWT:      true,
		userSDK:        userSDK,
		uf:             uf,
		UserContextKey: CurrentUserContextKey,
		JWTContextKey:  JwtContextKey,
	})
}

// CurrentUserWithoutFetch is for when you have a gateway that find the user and include
// it in the jwt. so you will dont need to any user finder.
func CurrentUserWithoutFetch(userSDK hexa.UserSDK) echo.MiddlewareFunc {
	return CurrentUserWithConfig(CurrentUserConfig{
		ExtendJWT:      false,
		uf:             nil,
		userSDK:        userSDK,
		UserContextKey: CurrentUserContextKey,
		JWTContextKey:  JwtContextKey,
	})
}

// CurrentUser is a middleware to set the user in the context.
// If provided jwt, so this function find user and set it as user
// otherwise set guest user.
func CurrentUserWithConfig(cfg CurrentUserConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {

			var user = cfg.userSDK.NewGuest()

			// Get jwt (if exists)
			if token, ok := ctx.Get(cfg.JWTContextKey).(*jwt.Token); ok {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					if cfg.ExtendJWT {
						user, err := cfg.uf(claims["sub"].(string))
						if err != nil {
							err = tracer.Trace(err)
							return err
						}
						extension, err := cfg.userSDK.Export(user)
						gutil.ExtendMap(claims, extension, true)
					}

					user, err = cfg.userSDK.Import(hexa.Map(claims))
					if err != nil {
						err = tracer.Trace(err)
						return
					}

				} else {
					return errors.New("JWT claims is not valid")
				}

			}

			// Set user in context with the given key
			ctx.Set(cfg.UserContextKey, user)

			// Also set for user to uas in hexa context
			ctx.Set(ContextKeyHexaUser, user)

			return next(ctx)
		}
	}
}
