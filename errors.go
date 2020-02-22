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
		true, http.StatusInternalServerError, "KE_1_E_0",
		kitty.ErrorKeyInternalError, "user not found",
	)

	errContextUserNotImplementedKittyUser = kitty.NewError(
		true, http.StatusInternalServerError, "KE_1_E_1",
		kitty.ErrorKeyInternalError, "context's user not implemented the kitty User interface.",
	)

	errJwtMissing = kitty.NewError(false, http.StatusBadRequest, "KE_1_E_2",
		"missing_jwt_token", "missing or malformed jwt")

	errInvalidOrExpiredJwt = kitty.NewError(false, http.StatusUnauthorized, "KE_1_E_3",
		"invalid_or_expired_jwt", "invalid or expired jwt")

	//--------------------------------
	// Request errors
	//--------------------------------
	errRequestIdNotFound = kitty.NewError(
		true, http.StatusInternalServerError, "KE_2_E_02",
		kitty.ErrorKeyInternalError, "request id not found in the request.",
	)

	//--------------------------------
	// Other errors
	//--------------------------------
	errUnknownError = kitty.NewError(true, http.StatusInternalServerError, "KE_3_E_00",
		"err_unknown_error","")
)
