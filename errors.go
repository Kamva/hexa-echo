package kecho

import (
	"github.com/Kamva/kitty"
	"net/http"
)

var (
	// Error code description:
	// KE = kitty echo
	// E = Error
	// 00 = error number.

	userNotFound = kitty.NewError(
		true, http.StatusInternalServerError, "KE_E_00",
		kitty.ErrorKeyInternalError, "user not found",
	)

	contextUserNotImplementedKittyUser = kitty.NewError(
		true, http.StatusInternalServerError, "KE_E_01",
		kitty.ErrorKeyInternalError, "context's user not implemented the kitty User interface.",
	)

	requestIdNotFound = kitty.NewError(
		true, http.StatusInternalServerError, "KE_E_02",
		kitty.ErrorKeyInternalError, "request id not found in the request.",
	)
)
