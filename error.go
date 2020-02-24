package kecho

import (
	"github.com/Kamva/kitty"
	"github.com/labstack/echo/v4"
)

// HTTPErrorHandler is the echo error handler.
// this function need to the KittyContext middleware.
func HTTPErrorHandler(requestErr error, c echo.Context) {
	if _, ok := requestErr.(kitty.Reply); !ok {
		requestErr = errUnknownError.SetInternalMessage(requestErr.Error())
	}

	kerr := requestErr.(kitty.Error)

	kittyCtx := c.Get(ContextKeyKittyCtx).(kitty.Context)
	translator := kittyCtx.Translator()
	msg, err := translator.Translate(kerr.Key(), kerr.Params())

	if err != nil {
		msg = ""
	}

	requestErr = c.JSON(kerr.HTTPStatus(), kitty.NewBody(kerr.Code(), msg, kitty.Data(kerr.Data())))

	// Report if need to report:
	if kerr.ShouldReport() {
		kerr.Report(kittyCtx.Logger(), kittyCtx.Translator())
	}

	if requestErr != nil {
		kittyCtx.Logger().Error(requestErr)
	}
}
