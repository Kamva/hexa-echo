package kecho

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const JwtContextKey="jwt"

// skipIfNotProvidedHeader skip jwt middleware if jwt authorization header
// is not provided.
func skipIfNotProvidedHeader(header string) middleware.Skipper {
	return func(c echo.Context) bool {
		return c.Request().Header.Get(header) != ""
	}
}

// jwtErrorHandler check errors type and return relative kitty error.
func jwtErrorHandler(err error) error {
	// missing or malformed jwt token
	if err == middleware.ErrJWTMissing {
		return errJwtMissing
	}

	// otherwise authorization error
	return errInvalidOrExpiredJwt
}

var jwtConfig = middleware.JWTConfig{
	Skipper:       skipIfNotProvidedHeader(echo.HeaderAuthorization),
	SigningMethod: middleware.AlgorithmHS256,
	ContextKey:    JwtContextKey,
	TokenLookup:   "header:" + echo.HeaderAuthorization,
	AuthScheme:    "Bearer",
	Claims:        jwt.MapClaims{},
	ErrorHandler:  jwtErrorHandler,
}

// JWT middleware
func JWT(key string) echo.MiddlewareFunc {
	cfg := jwtConfig
	cfg.SigningKey = key

	return middleware.JWTWithConfig(cfg)
}
