package kecho

import (
	"github.com/Kamva/kitty"
	"github.com/Kamva/tracer"
	"github.com/labstack/echo/v4"
)

const (
	// ContextKeyKittyRequestID uses as key in context to store request id to use in context middleware
	ContextKeyKittyRequestID = "__kitty_ctx.rid__"

	// ContextKeyKittyCorrelationID uses as key in context to store correlation id to use in context middleware
	ContextKeyKittyCorrelationID = "__kitty_ctx.cid__"

	// ContextKeyKittyCtx is the identifier to set the kitty context as a field in the context of a request.
	// e.g ctx.Set(kitty.ContextIdentifier,kittyCtx) // kittyCtx is kitty Context.
	ContextKeyKittyCtx = "__kitty_ctx.ctx__"

	// ContextKeyKittyUser is the identifier to set the kitty user as a field in the context of a request.
	ContextKeyKittyUser = "__kitty_ctx.user__"
)

// getKittyUser returns kitty user instance from the current user.
func getKittyUser(ctx echo.Context) (kitty.User, kitty.Error) {
	// Get user if exists:
	u := ctx.Get(ContextKeyKittyUser)

	if u == nil {
		return nil, errUserNotFound
	}

	if u, ok := u.(kitty.User); ok {
		return u, nil
	} else {
		return nil, errContextUserNotImplementedKittyUser
	}
}

// getRequestID returns the request id.
func getRequestID(ctx echo.Context) (string, kitty.Error) {
	// Get Request ID if exists:
	rid := ctx.Get(ContextKeyKittyRequestID).(string)

	if rid == "" {
		return "", errRequestIdNotFound
	}

	return rid, nil
}

// getCorrelationID returns the request correlation id.
func getCorrelationID(ctx echo.Context) (string, kitty.Error) {
	// Get Request ID if exists:
	cid := ctx.Get(ContextKeyKittyCorrelationID).(string)

	if cid == "" {
		return "", errCorrelationIDNotFound
	}

	return cid, nil
}

// KittyContext set kitty context on each request.
func KittyContext(logger kitty.Logger, translator kitty.Translator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			req := ctx.Request()

			user, err := getKittyUser(ctx)

			if err != nil {
				return tracer.Trace(err)
			}

			rid, err := getRequestID(ctx)

			if err != nil {
				return tracer.Trace(err)
			}

			cid, err := getCorrelationID(ctx)

			if err != nil {
				return tracer.Trace(err)
			}

			al := req.Header.Get("Accept-Language")

			// Set context
			ctx.Set(ContextKeyKittyCtx, kitty.NewCtx(rid, cid, al, user, logger, translator))

			// Set context logger
			ctx.SetLogger(KittyLoggerToEchoLogger(logger))

			return next(ctx)
		}
	}
}