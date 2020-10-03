package hecho

import (
	"errors"
	"fmt"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"github.com/dgrijalva/jwt-go"
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
		Secret                  hexa.Secret
		ExpireTokenAfter        time.Duration
		ExpireRefreshTokenAfter time.Duration
		SubGenerator            SubGenerator
	}

	// RefreshTokenConfig use as config to refresh access token.
	RefreshTokenConfig struct {
		GenerateTokenConfig
		RefreshToken hexa.Secret
		// Use Authorizer to verify that can get new token.
		Authorizer RefreshTokenAuthorizer
	}
)

//--------------------------------
// JWT Middleware
//--------------------------------

const JwtContextKey = "jwt"

// skipIfNotProvidedHeader skip jwt middleware if jwt authorization header
// is not provided.
func skipIfNotProvidedHeader(header string) middleware.Skipper {
	return func(c echo.Context) bool {
		return c.Request().Header.Get(header) == ""
	}
}

// jwtErrorHandler check errors type and return relative hexa error.
func jwtErrorHandler(err error) error {
	// missing or malformed jwt token
	if err == middleware.ErrJWTMissing {
		return errJwtMissing.SetError(tracer.Trace(err))
	}

	// otherwise authorization error
	return errInvalidOrExpiredJwt.SetError(tracer.Trace(err))
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
func JWT(key hexa.Secret) echo.MiddlewareFunc {
	cfg := jwtConfig
	cfg.SigningKey = []byte(key)
	// TODO: remove this function and user config, provide that config instead of secret in huner, also provide config for
	// cont.TODO: generate and refresh token, change RefreshToken function to check for RSA algorithm also.
	return middleware.JWTWithConfig(cfg)
}

//--------------------------------
// JWT Generator
//--------------------------------

// IDAsSubjectGenerator return user's id as jwt subject (sub).
func IDAsSubjectGenerator(user hexa.User) (string, error) {
	return user.Identifier().String(), nil
}

// GenerateToken generate new token for the user.
func GenerateToken(u hexa.User, cfg GenerateTokenConfig) (token, rToken hexa.Secret, err error) {
	if err = tracer.Trace(validateGenerateTokenCfg(cfg)); err != nil {
		return
	}

	sub, err := cfg.SubGenerator(u)
	if err != nil {
		err = tracer.Trace(err)
		return
	}

	// Generate Token
	token, err = generateToken(jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(cfg.ExpireTokenAfter).Unix(),
	}, cfg.Secret)

	if err != nil {
		err = tracer.Trace(err)
		return
	}

	// Generate Refresh token
	rToken, err = generateToken(jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(cfg.ExpireRefreshTokenAfter).Unix(),
	}, cfg.Secret)
	err = tracer.Trace(err)
	return
}

// RefreshToken refresh the jwt token by provided config.
// In provided config to this function set user as just simple
// hexa guest user. we set it by your authorizer later.
func RefreshToken(cfg RefreshTokenConfig) (token, rToken hexa.Secret, err error) {
	if err = tracer.Trace(validateRefreshTokenCfg(cfg)); err != nil {
		return
	}

	// Parse token:
	jToken, err := jwt.Parse(string(cfg.RefreshToken), func(token *jwt.Token) (interface{}, error) {
		// validate hashing method.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, tracer.Trace(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}

		return []byte(cfg.Secret), nil
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
	if cfg.ExpireTokenAfter > cfg.ExpireRefreshTokenAfter {
		return errors.New("refresh token expire time can not be less than access token expire time")

	}

	if cfg.Secret == "" {
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

// generateToken generate new jwt token.
func generateToken(claims jwt.MapClaims, secret hexa.Secret) (token hexa.Secret, err error) {
	jToken := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	jToken.Claims = claims
	// Generate encoded token and send it as response.
	t, err := jToken.SignedString([]byte(secret))
	if err != nil {
		err = tracer.Trace(err)
		return
	}
	token = hexa.Secret(t)
	return
}
