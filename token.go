package hecho

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const TokenHeaderAuthorization = "Authorization"
const TokenCookieFieldAuthToken = "hexa_auth_token"

const AuthTokenContextKey = "auth_token"
const AuthTokenLocationContextKey = "auth_token_location"

type TokenLocation int

const (
	TokenLocationCookie = iota
	TokenLocationHeader
)

type ExtractTokenConfig struct {
	Skipper                 middleware.Skipper
	TokenHeaderField        string
	TokenCookieField        string
	TokenContextKey         string
	TokenLocationContextKey string
}

func ExtractAuthToken() echo.MiddlewareFunc {
	return ExtractTokenWithConfig(ExtractTokenConfig{
		TokenHeaderField:        TokenHeaderAuthorization,
		TokenCookieField:        TokenCookieFieldAuthToken,
		TokenContextKey:         AuthTokenContextKey,
		TokenLocationContextKey: AuthTokenLocationContextKey,
	})
}

// ExtractTokenWithConfig extracts the authentication token from the cookie or Authorization header.
func ExtractTokenWithConfig(cfg ExtractTokenConfig) echo.MiddlewareFunc {
	if cfg.Skipper == nil {
		cfg.Skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if cfg.Skipper(ctx) {
				return next(ctx)
			}

			// Fetch from cookies
			if cookie, err := ctx.Cookie(cfg.TokenCookieField); err == nil {
				ctx.Set(AuthTokenContextKey, cookie.Value)
				ctx.Set(AuthTokenLocationContextKey, TokenLocationCookie)
				return next(ctx)
			}

			// Fetch from header
			if headerVal := ctx.Request().Header.Get(cfg.TokenHeaderField); headerVal != "" {
				var token string
				if len(headerVal) > 6 && strings.ToUpper(headerVal[0:6]) == "BEARER" {
					token = headerVal[7:]
				}

				if len(headerVal) > 5 && strings.ToUpper(headerVal[0:5]) == "TOKEN" {
					token = headerVal[6:]
				}

				if token != "" {
					ctx.Set(cfg.TokenContextKey, token)
					ctx.Set(cfg.TokenLocationContextKey, TokenLocationHeader)
					return next(ctx)
				}
			}

			return next(ctx)
		}
	}
}
