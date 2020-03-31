package hecho

import (
	"errors"
	"github.com/Kamva/hexa"
	"net/http"
)

var (
	// Error code description:
	// hec = hexa echo (package or project name)
	// 1 = errors about user section (identify some part in application)
	// E = Error (type of code : error|response)
	// 00 = error number zero (id of code in that part and type)

	//--------------------------------
	// User and authentication Errors
	//--------------------------------
	errUserNotFound = hexa.NewError(http.StatusInternalServerError, "hec.ua.e.0",
		hexa.ErrKeyInternalError, errors.New("user not found"),
	)

	errContextUserNotImplementedHexaUser = hexa.NewError(http.StatusInternalServerError, "hec.ua.e.1",
		hexa.ErrKeyInternalError, errors.New("context's user not implemented the hexa User interface"),
	)

	errJwtMissing = hexa.NewError(http.StatusBadRequest, "hec.ua.e.2",
		"err_missing_jwt_token", errors.New("missing or malformed jwt"))

	errInvalidOrExpiredJwt = hexa.NewError(http.StatusUnauthorized, "hec.ua.e.3",
		"err_invalid_or_expired_jwt", errors.New("invalid or expired jwt"))

	errUserMustBeGuest = hexa.NewError(http.StatusUnauthorized, "hec.ua.e.4",
		"err_user_must_be_guest", errors.New("the user must be guest to access this API"))

	errUserNeedToAuthenticate = hexa.NewError(http.StatusUnauthorized, "hec.ua.e.5",
		"err_user_must_authenticate", errors.New("the user need to login to access this API"))

	errRefreshTokenCanNotBeEmpty = hexa.NewError(http.StatusBadRequest, "hec.ua.e.6",
		"err_refresh_token_is_empty", errors.New("refresh token can not be empty"))

	errInvalidRefreshToken = hexa.NewError(http.StatusBadRequest, "hec.ua.e.7",
		"err_invalid_refresh_token", nil)

	//--------------------------------
	// Request errors
	//--------------------------------
	errRequestIdNotFound = hexa.NewError(http.StatusInternalServerError, "hec.rq.e.2",
		hexa.ErrKeyInternalError, errors.New("request id not found in the request"),
	)

	errCorrelationIDNotFound = hexa.NewError(http.StatusInternalServerError, "hec.rq.e.3",
		hexa.ErrKeyInternalError, errors.New("correlation id not found in the request"),
	)

	//--------------------------------
	// DEBUG
	//--------------------------------
	errRouteAvaialbeInDebugMode = hexa.NewError(http.StatusUnauthorized, "hec.dbg.e.0",
		"err_route_available_in_debug_mode", errors.New("route is available just in debug mode"),
	)

	//--------------------------------
	// Other errors
	//--------------------------------
	errHTTPNotFoundError = hexa.NewError(http.StatusNotFound, "hec.ot.e.0", "route_not_found", nil)

	// Set this error status manually on return relative to echo error code.
	errEchoHTTPError = hexa.NewError(http.StatusNotFound, "hec.ot.e.1", hexa.TranslateKeyEmptyMessage, nil)

	errUnknownError = hexa.NewError(http.StatusInternalServerError, "hec.ot.e.1", "err_unknown_error", nil)
)
