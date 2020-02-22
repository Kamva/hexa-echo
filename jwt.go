package kecho

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strings"
)

// skipIfNotProvidedHeader skip jwt middleware if jwt authorization header
// is not provided.
func skipIfNotProvidedHeader(ctx echo.Context) bool {
	parts := strings.Split(jwtConfig.TokenLookup, ":")
	return ctx.Request().Header.Get(parts[1]) != ""
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
	Skipper:       skipIfNotProvidedHeader,
	SigningMethod: middleware.AlgorithmHS256,
	ContextKey:    "jwt",
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
