package hecho

import (
	"errors"
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
	"github.com/Kamva/tracer"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type (
	// JWTExtender is a function to get the sub and return map of
	// extension values for jwt.
	// Our pattern in authentication has this instruction:
	// - Get jwt token (later can get api,oauth2,... tokens). here we just get jwt token, but later in API Gateway you
	//  can get other types of tokens and extend jwt, then send extended jwt to microservices and disable
	// "ExtendJWT" to true to prevent double extension of jwt(one time in gateway, another time in microservice).
	// - Call to extender(if needed) to extend jwt token by adding some values to it (e.g permissionsList,...).
	// - Call to user generator to generate user by jwt token.
	JWTExtender func(sub string) (map[string]interface{}, error)

	// UserGeneratorByJWT generate new user by extended jwt token.
	// This function get jwt claims as its first argument.
	UserGeneratorByExtendedJWT func(claims map[string]interface{}) (hexa.User, error)

	// CurrentUserConfig is the config to use in CurrentUser middleware.
	CurrentUserConfig struct {
		JWTExtender
		UserGeneratorByExtendedJWT
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
func CurrentUser(jwtExtender JWTExtender, userGenerator UserGeneratorByExtendedJWT) echo.MiddlewareFunc {
	return CurrentUserWithConfig(CurrentUserConfig{
		ExtendJWT:                  true,
		JWTExtender:                jwtExtender,
		UserGeneratorByExtendedJWT: userGenerator,
		UserContextKey:             CurrentUserContextKey,
		JWTContextKey:              JwtContextKey,
	})
}

func CurrentUserWithoutExtender(userGenerator UserGeneratorByExtendedJWT) echo.MiddlewareFunc {
	return CurrentUserWithConfig(CurrentUserConfig{
		ExtendJWT:                  false,
		JWTExtender:                nil,
		UserGeneratorByExtendedJWT: userGenerator,
		UserContextKey:             CurrentUserContextKey,
		JWTContextKey:              JwtContextKey,
	})
}

// CurrentUser is a middleware to set the user in the context.
// If provided jwt, so this function find user and set it as user
// otherwise set guest user.
func CurrentUserWithConfig(cfg CurrentUserConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {

			var user = hexa.NewGuestUser()

			// Get jwt (if exists)
			if token, ok := ctx.Get(cfg.JWTContextKey).(*jwt.Token); ok {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					if cfg.ExtendJWT {
						extension, err := cfg.JWTExtender(claims["sub"].(string))
						if err != nil {
							err = tracer.Trace(err)
							return err
						}
						gutil.ExtendMap(claims, extension, true)
					}

					user, err = cfg.UserGeneratorByExtendedJWT(claims)
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