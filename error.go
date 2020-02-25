package kecho

import (
	"github.com/Kamva/kitty"
	"github.com/labstack/echo/v4"
)

// HTTPErrorHandler is the echo error handler.
// this function need to the KittyContext middleware.
func HTTPErrorHandler(l kitty.Logger, t kitty.Translator) echo.HTTPErrorHandler {
	return func(requestErr error, c echo.Context) {
		l := l
		t := t

		if _, ok := requestErr.(kitty.Reply); !ok {
			requestErr = errUnknownError.SetInternalMessage(requestErr.Error())
		}

		kerr := requestErr.(kitty.Error)

		kittyCtx := c.Get(ContextKeyKittyCtx).(kitty.Context)

		// Maybe error occur before set kitty context in middleware
		if kittyCtx != nil {
			l = kittyCtx.Logger()
			t = kittyCtx.Translator()
		}

		msg, err := t.Translate(kerr.Key(), kerr.Params())

		if err != nil {
			msg = ""
		}

		requestErr = c.JSON(kerr.HTTPStatus(), kitty.NewBody(kerr.Code(), msg, kitty.Data(kerr.Data())))

		// Report if need to report:
		if kerr.ShouldReport() {
			kerr.Report(l, t)
		}

		if requestErr != nil {
			l.Error(requestErr)
		}
	}

}
