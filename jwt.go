package kecho

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Kamva/kitty"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

type (

	// RefreshTokenAuthorizer is a type check that user can get new token.
	RefreshTokenAuthorizer func(sub string) (kitty.User, error)

	SubGenerator func(user kitty.User) (string, error)

	// GenerateTokenConfig use as config to generate new token.
	GenerateTokenConfig struct {
		User                    kitty.User
		Secret                  kitty.Secret
		ExpireTokenAfter        time.Duration
		ExpireRefreshTokenAfter time.Duration
		SubGenerator            SubGenerator
	}

	// RefreshTokenConfig use as config to refresh access token.
	RefreshTokenConfig struct {
		GenerateTokenConfig
		RefreshToken kitty.Secret
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

//--------------------------------
// JWT Generator
//--------------------------------

// IDAsSubjectGenerator return user's id as jwt subject (sub).
func IDAsSubjectGenerator(user kitty.User) (string, error) {
	id := user.GetID()
	switch id.(type) {
	case string:
		return id.(string), nil
	case json.Marshaler:
		m, err := id.(json.Marshaler).MarshalJSON()
		return string(m), err
	}

	return fmt.Sprint(id), nil
}

// GenerateToken generate new token for the user.
func GenerateToken(cfg GenerateTokenConfig) (token, rToken kitty.Secret, err error) {
	if err = validateGenerateTokenCfg(cfg); err != nil {
		return
	}

	sub, err := cfg.SubGenerator(cfg.User)
	if err != nil {
		return
	}

	// Generate Token
	token, err = generateToken(jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(cfg.ExpireTokenAfter).Unix(),
	}, cfg.Secret)

	if err != nil {
		return
	}

	// Generate Refresh token
	rToken, err = generateToken(jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(cfg.ExpireRefreshTokenAfter).Unix(),
	}, cfg.Secret)

	return
}

// RefreshToken refresh the jwt token by provided config.
// In provided config to this function set user as just simple
// kitty guest user. we set it by your authorizer later.
func RefreshToken(cfg RefreshTokenConfig) (token, rToken kitty.Secret, err error) {
	if err = validateRefreshTokenCfg(cfg); err != nil {
		return
	}

	// Parse token:
	jToken, err := jwt.Parse(string(cfg.RefreshToken), func(token *jwt.Token) (interface{}, error) {
		// validate hashing method.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return
	}

	// Authorize user to verify user can get new access token.
	user, err := cfg.Authorizer(jToken.Claims.(jwt.MapClaims)["sub"].(string))

	if err != nil {
		return
	}

	// Set provided user.
	cfg.GenerateTokenConfig.User = user

	return GenerateToken(cfg.GenerateTokenConfig)
}

func validateGenerateTokenCfg(cfg GenerateTokenConfig) error {
	if cfg.ExpireTokenAfter > cfg.ExpireRefreshTokenAfter {
		return errors.New("refresh token expire time can not be less than access token expire time")

	}

	if cfg.User == nil || cfg.Secret == "" {
		return errors.New("invalid config values to generate token pairs")
	}

	if cfg.SubGenerator == nil {
		return errors.New("invalid subject generator for jwt")
	}

	return nil
}

func validateRefreshTokenCfg(cfg RefreshTokenConfig) error {
	if err := validateGenerateTokenCfg(cfg.GenerateTokenConfig); err != nil {
		return err
	}

	if cfg.Authorizer == nil {
		return errors.New("authorizer can not be nil")
	}

	if cfg.RefreshToken == "" {
		return errors.New("refresh token can not be empty")
	}

	return nil
}

// generateToken generate new jwt token.
func generateToken(claims jwt.MapClaims, secret kitty.Secret) (token kitty.Secret, err error) {
	jToken := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	jToken.Claims = claims

	// Generate encoded token and send it as response.
	t, err := jToken.SignedString([]byte(secret))

	if err != nil {
		return
	}

	token = kitty.Secret(t)

	return
}
