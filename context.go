package hecho

import (
	"github.com/Kamva/hexa"
	"github.com/Kamva/tracer"
	"github.com/labstack/echo/v4"
)

const (
	// ContextKeyHexaRequestID uses as key in context to store request id to use in context middleware
	ContextKeyHexaRequestID = "__hexa_ctx.rid__"

	// ContextKeyHexaCorrelationID uses as key in context to store correlation id to use in context middleware
	ContextKeyHexaCorrelationID = "__hexa_ctx.cid__"

	// ContextKeyHexaCtx is the identifier to set the hexa context as a field in the context of a request.
	ContextKeyHexaCtx = "__hexa_ctx.ctx__"

	// ContextKeyHexaUser is the identifier to set the hexa user as a field in the context of a request.
	ContextKeyHexaUser = "__hexa_ctx.user__"
)

// getHexaUser returns hexa user instance from the current user.
func getHexaUser(ctx echo.Context) (hexa.User, hexa.Error) {
	// Get user if exists:
	u := ctx.Get(ContextKeyHexaUser)

	if u == nil {
		return nil, errUserNotFound
	}

	if u, ok := u.(hexa.User); ok {
		return u, nil
	} else {
		return nil, errContextUserNotImplementedHexaUser
	}
}

// getCorrelationID returns the request correlation id.
func getCorrelationID(ctx echo.Context) (string, hexa.Error) {
	// Get Request ID if exists:
	cid := ctx.Get(ContextKeyHexaCorrelationID).(string)

	if cid == "" {
		return "", errCorrelationIDNotFound
	}

	return cid, nil
}

// HexaContext set hexa context on each request.
func HexaContext(logger hexa.Logger, translator hexa.Translator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			req := ctx.Request()

			user, err := getHexaUser(ctx)

			if err != nil {
				return tracer.Trace(err)
			}

			cid, err := getCorrelationID(ctx)

			if err != nil {
				return tracer.Trace(err)
			}

			al := req.Header.Get("Accept-Language")

			// Set context
			ctx.Set(ContextKeyHexaCtx, hexa.NewCtx(ctx.Request(), cid, al, user, logger, translator))

			return next(ctx)
		}
	}
}

