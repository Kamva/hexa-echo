package hecho

import (
	"math"
	"net/http"
	"strconv"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/throttled/throttled/v2"
)

var ErrTooManyRequests = hexa.NewError(http.StatusTooManyRequests, "lib.http.too_many_requests_error", nil)

const (
	// HeaderRateLimitLimit, HeaderRateLimitRemaining, and HeaderRateLimitReset
	// are the recommended return header values from IETF on rate limiting. Reset
	// is in UTC time.

	HeaderRateLimitLimit     = "X-RateLimit-Limit"
	HeaderRateLimitRemaining = "X-RateLimit-Remaining"
	HeaderRateLimitReset     = "X-RateLimit-Reset"

	// HeaderRetryAfter is the header used to indicate when a client should retry
	// requests (when the rate limit expires), in UTC time.
	HeaderRetryAfter = "Retry-After"
)

type KeyExtractor func(c echo.Context) (string, error)

type RateLimiterConfig struct {
	Skipper      middleware.Skipper
	RateLimiter  throttled.RateLimiter
	KeyExtractor KeyExtractor
}

func RateLimiter(cfg RateLimiterConfig) echo.MiddlewareFunc {
	if cfg.Skipper == nil {
		cfg.Skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			k, err := cfg.KeyExtractor(c)
			if err != nil {
				return tracer.Trace(err)
			}

			limited, context, err := cfg.RateLimiter.RateLimit(k, 1)

			if err != nil {
				hlog.Error("error on checking rate limit", hlog.Err(tracer.Trace(err)))
				return tracer.Trace(err)
			}

			resetAfterSeconds := int(math.Ceil(context.ResetAfter.Seconds()))
			c.Response().Header().Set(HeaderRateLimitLimit, strconv.FormatInt(int64(context.Limit), 10))
			c.Response().Header().Set(HeaderRateLimitRemaining, strconv.FormatInt(int64(context.Remaining), 10))
			c.Response().Header().Set(HeaderRateLimitReset, strconv.FormatInt(int64(resetAfterSeconds), 10))

			// Fail if there were no tokens remaining.

			if limited {
				retryAfterSeconds := int(math.Ceil(context.RetryAfter.Seconds()))
				c.Response().Header().Set(HeaderRetryAfter, strconv.FormatInt(int64(retryAfterSeconds), 10))
				return ErrTooManyRequests
			}

			return next(c)
		}
	}
}

func RealIPKeyExtractor(c echo.Context) (string, error) {
	return c.RealIP(), nil
}
