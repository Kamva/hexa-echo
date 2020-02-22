package kecho

import (
	"github.com/Kamva/elogrus/v4"
	"github.com/Kamva/kitty"
	"github.com/Kamva/kitty/kittylogger"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const (
	userNotFound      = "user not found in the context"
	requestIdNotFound = "request id not found in the request."
)

const KittyLoggerKey = "kitty.logger"

// ContextLogSetter set custom logger for each context.
func ContextLogSetter(logger *logrus.Entry) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// Get user if exists:
			if u, ok := ctx.Get("user").(kitty.User); ok {
				logger = logger.WithFields(logrus.Fields{
					"guest":    u.IsGuest(),
					"user_id":  u.GetID(),
					"username": u.GetUsername(),
				})
			} else {
				logger.WithFields(logrus.Fields{"from": "context log setter"}).Error(userNotFound)
			}

			req := ctx.Request()
			// Get Request ID if exists:
			rid := req.Header.Get(echo.HeaderXRequestID)

			if rid != "" {
				logger = logger.WithFields(logrus.Fields{"request_id": rid})
			} else {
				logger.WithFields(logrus.Fields{"from": "context log setter"}).Error(requestIdNotFound)
			}

			// Set logger as context standard logger and kitty logger in ctx.
			ctx.SetLogger(elogrus.GetEchoLogger(logger))

			// Set kitty logger
			ctx.Set(KittyLoggerKey, kittylogger.NewLogrusDriver(logger))

			return next(ctx)
		}
	}
}
