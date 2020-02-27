package kecho

import (
	"errors"
	"github.com/Kamva/kitty"
	"net/http"
)

var (
	// Error code description:
	// KE = kitty echo (package or project name)
	// 1 = errors about user section (identify some part in application)
	// E = Error (type of code : error|response)
	// 00 = error number zero (id of code in that part and type)

	//--------------------------------
	// User and authentication Errors
	//--------------------------------
	errUserNotFound = kitty.NewError(http.StatusInternalServerError, "ke.1.e.0",
		kitty.ErrKeyInternalError, errors.New("user not found"),
	)

	errContextUserNotImplementedKittyUser = kitty.NewError(http.StatusInternalServerError, "ke.1.e.1",
		kitty.ErrKeyInternalError, errors.New("context's user not implemented the kitty User interface"),
	)

	errJwtMissing = kitty.NewError(http.StatusBadRequest, "ke.1.e.2",
		"missing_jwt_token", errors.New("missing or malformed jwt"))

	errInvalidOrExpiredJwt = kitty.NewError(http.StatusUnauthorized, "ke.1.e.3",
		"invalid_or_expired_jwt", errors.New("invalid or expired jwt"))

	errUserMustBeGuest = kitty.NewError(http.StatusUnauthorized, "ke.1.e.4",
		"user_must_be_guest", errors.New("the user must be guest to access this API"))

	errUserNeedToAuthenticate = kitty.NewError(http.StatusUnauthorized, "ke.1.e.5",
		"user_must_authenticate", errors.New("the user need to login to access this API"))

	//--------------------------------
	// Request errors
	//--------------------------------
	errRequestIdNotFound = kitty.NewError(http.StatusInternalServerError, "ke.2.e.2",
		kitty.ErrKeyInternalError, errors.New("request id not found in the request"),
	)

	errCorrelationIDNotFound = kitty.NewError(http.StatusInternalServerError, "ke.2.e.3",
		kitty.ErrKeyInternalError, errors.New("correlation id not found in the request"),
	)

	//--------------------------------
	// DEBUG
	//--------------------------------
	errRouteAvaialbeInDebugMode = kitty.NewError(http.StatusUnauthorized, "ke.3.e.0",
		"route_available_in_debug_mode", errors.New("route is available just in debug mode"),
	)

	//--------------------------------
	// Other errors
	//--------------------------------
	// Set this error status manually on return relative to echo error code.
	errEchoHTTPError = kitty.NewError(http.StatusNotFound, "ke.4.e.0", kitty.TranslateKeyEmptyMessage, nil)

	errUnknownError = kitty.NewError(http.StatusInternalServerError, "ke.4.e.1", "err_unknown_error", nil)
)
