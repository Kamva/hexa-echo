package kecho

import (
	"github.com/Kamva/kitty"
	"github.com/labstack/echo/v4"
)

// ContextKeyKittyCtx is the identifier to set the kitty context as a field in the context of a request.
// e.g ctx.Set(kitty.ContextIdentifier,kittyCtx) // kittyCtx is kitty Context.
const ContextKeyKittyCtx = "KITTY_CONTEXT"

// ContextKeyKittyUser is the identifier to set the kitty user as a field in the context of a request.
const ContextKeyKittyUser = "KITTY_USER"

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
	req := ctx.Request()
	// Get Request ID if exists:
	rid := req.Header.Get(echo.HeaderXRequestID)

	if rid == "" {
		return "", errRequestIdNotFound
	}

	return rid, nil
}

// getCorrelationID returns the request correlation id.
func getCorrelationID(ctx echo.Context) (string, kitty.Error) {
	req := ctx.Request()
	// Get Request ID if exists:
	rid := req.Header.Get(HeaderCorrelationID)

	if rid == "" {
		return "", errCorrelationIDNotFound
	}

	return rid, nil
}

// tuneLogger function tune the logger for users request.
func tuneLogger(ctx echo.Context, requestID string, correlationID string, u kitty.User, logger kitty.Logger) kitty.Logger {

	logger = logger.WithFields(
		"__guest__", u.IsGuest(),
		"__user_id__", u.GetID(),
		"__username__", u.GetUsername(),
		"__request_id__", requestID,
		"__correlation_id__", correlationID,
	)

	return logger
}

// localizeTranslator localize translator for each request relative to the request headers.
func localizeTranslator(ctx echo.Context, t kitty.Translator) kitty.Translator {
	req := ctx.Request()
	// Get Request ID if exists:
	al := req.Header.Get("Accept-Language")
	if al != "" {
		return t.Localize(al)
	}

	return t.Localize()
}

// KittyContext set kitty context on each request.
func KittyContext(logger kitty.Logger, translator kitty.Translator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			user, err := getKittyUser(ctx)

			if err != nil {
				return err
			}

			rid, err := getRequestID(ctx)

			if err != nil {
				return err
			}

			cid, err := getCorrelationID(ctx)

			if err != nil {
				return err
			}

			logger := tuneLogger(ctx, rid, cid, user, logger)

			translator := localizeTranslator(ctx, translator)

			// Set context
			ctx.Set(ContextKeyKittyCtx, kitty.NewCtx(rid, cid, user, logger, translator))

			// Set context logger
			ctx.SetLogger(KittyLoggerToEchoLogger(logger))

			return next(ctx)
		}
	}
}
