package hecho

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

type (

	// RefreshTokenAuthorizer is a type check that user can get new token.
	RefreshTokenAuthorizer func(sub string) (hexa.User, error)

	SubGenerator func(user hexa.User) (string, error)

	// GenerateTokenConfig use as config to generate new token.
	GenerateTokenConfig struct {
		SingingMethod    jwt.SigningMethod
		Key              interface{}
		SubGenerator     SubGenerator
		Claims           jwt.MapClaims
		ExpireTokenAfter time.Duration
	}

	// RefreshTokenConfig use as config to refresh access token.
	RefreshTokenConfig struct {
		GenerateTokenConfig
		RefreshToken string
		// Use Authorizer to verify that can get new token.
		Authorizer RefreshTokenAuthorizer
	}
)

//--------------------------------
// JWT Middleware
//--------------------------------

const JwtContextKey = "jwt"

// SkipIfNotProvidedHeader skip jwt middleware if jwt authorization header
// is not provided.
func SkipIfNotProvidedHeader(header string) middleware.Skipper {
	return func(c echo.Context) bool {
		return c.Request().Header.Get(header) == ""
	}
}

// JwtErrorHandler check errors type and return relative hexa error.
func JwtErrorHandler(err error) error {
	// missing or malformed jwt token
	if err == middleware.ErrJWTMissing {
		return errJwtMissing.SetError(tracer.Trace(err))
	}

	// otherwise authorization error
	return errInvalidOrExpiredJwt.SetError(tracer.Trace(err))
}

//--------------------------------
// JWT Generator
//--------------------------------

// IDAsSubjectGenerator return user's id as jwt subject (sub).
func IDAsSubjectGenerator(user hexa.User) (string, error) {
	return user.Identifier().String(), nil
}

// GenerateToken generate new token for the user.
func GenerateToken(u hexa.User, cfg GenerateTokenConfig) (token string, err error) {
	if err = tracer.Trace(validateGenerateTokenCfg(cfg)); err != nil {
		return
	}

	sub, err := cfg.SubGenerator(u)
	if err != nil {
		err = tracer.Trace(err)
		return
	}
	gutil.ExtendMap(cfg.Claims, jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(cfg.ExpireTokenAfter).Unix(),
	}, true)

	jToken := jwt.New(cfg.SingingMethod)
	// Set claims
	jToken.Claims = cfg.Claims
	// Generate encoded token and send it as response.
	t, err := jToken.SignedString(cfg.Key)
	return t, tracer.Trace(err)
}

// GenerateRefreshToken refresh the jwt token by provided config.
func GenerateRefreshToken(cfg RefreshTokenConfig) (token string, err error) {
	if err = tracer.Trace(validateRefreshTokenCfg(cfg)); err != nil {
		return
	}

	// Parse token:
	jToken, err := jwt.Parse(cfg.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return cfg.Key, nil
	})

	if err != nil {
		err = errInvalidRefreshToken.SetError(tracer.Trace(err))
		return
	}

	// Authorize user to verify user can get new access token.
	user, err := cfg.Authorizer(jToken.Claims.(jwt.MapClaims)["sub"].(string))

	if err != nil {
		err = tracer.Trace(err)
		return
	}

	return GenerateToken(user, cfg.GenerateTokenConfig)
}

func validateGenerateTokenCfg(cfg GenerateTokenConfig) error {
	if gutil.IsNil(cfg.SingingMethod) || gutil.IsNil(cfg.Key) {
		return tracer.Trace(errors.New("invalid config values to generate token pairs"))
	}

	if cfg.SubGenerator == nil {
		return tracer.Trace(errors.New("invalid subject uuidGenerator for jwt"))
	}

	return nil
}

func validateRefreshTokenCfg(cfg RefreshTokenConfig) error {
	if err := validateGenerateTokenCfg(cfg.GenerateTokenConfig); err != nil {
		return tracer.Trace(err)
	}

	if cfg.Authorizer == nil {
		return tracer.Trace(errors.New("authorizer can not be nil"))
	}

	if cfg.RefreshToken == "" {
		return errRefreshTokenCanNotBeEmpty
	}

	return nil
}
