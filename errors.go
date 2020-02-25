package kecho

import (
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
	errUserNotFound = kitty.NewError(
		true, http.StatusInternalServerError, "ke.1.e.0",
		kitty.ErrorKeyInternalError, "user not found.",
	)

	errContextUserNotImplementedKittyUser = kitty.NewError(
		true, http.StatusInternalServerError, "ke.1.e.1",
		kitty.ErrorKeyInternalError, "context's user not implemented the kitty User interface.",
	)

	errJwtMissing = kitty.NewError(false, http.StatusBadRequest, "ke.1.e.2",
		"missing_jwt_token", "missing or malformed jwt.")

	errInvalidOrExpiredJwt = kitty.NewError(false, http.StatusUnauthorized, "ke.1.e.3",
		"invalid_or_expired_jwt", "invalid or expired jwt.")

	errUserMustBeGuest = kitty.NewError(false, http.StatusUnauthorized, "ke.1.e.4",
		"user_must_be_guest", "The user must be guest to access this API.")

	errUserNeedToAuthenticate = kitty.NewError(false, http.StatusUnauthorized, "ke.1.e.5",
		"user_must_authenticate", "The user need to login to access this API.")

	//--------------------------------
	// Request errors
	//--------------------------------
	errRequestIdNotFound = kitty.NewError(
		true, http.StatusInternalServerError, "ke.2.e.2",
		kitty.ErrorKeyInternalError, "request id not found in the request.",
	)

	errCorrelationIDNotFound = kitty.NewError(
		true, http.StatusInternalServerError, "ke.2.e.3",
		kitty.ErrorKeyInternalError, "correlation id not found in the request.",
	)

	//--------------------------------
	// Other errors
	//--------------------------------
	errUnknownError = kitty.NewError(true, http.StatusInternalServerError, "ke.3.e.0",
		"err_unknown_error", "")
)
