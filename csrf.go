package hecho

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CSRFSkipper skips if request doesn't need to csrf check.
// We do csrf check when user's token is in the cookie or session.
// and the request method is post too.
func CSRFSkipper(ctx echo.Context) bool {
	l := ctx.Get(AuthTokenLocationContextKey)
	return l != TokenLocationCookie && l != TokenLocationSession
}

var _ middleware.Skipper = CSRFSkipper
