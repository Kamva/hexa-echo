package kecho

import (
	"github.com/Kamva/gutil"
	"github.com/Kamva/kitty"
	"github.com/labstack/echo/v4"
)

// HTTPErrorHandler is the echo error handler.
// this function need to the KittyContext middleware.
func HTTPErrorHandler(l kitty.Logger, t kitty.Translator) echo.HTTPErrorHandler {
	return func(requestErr error, c echo.Context) {
		l := l
		t := t

		if httpErr, ok := requestErr.(*echo.HTTPError); ok {

			requestErr = errEchoHTTPError.SetShouldReport(httpErr.Code >= 500)
			requestErr = errEchoHTTPError.SetHTTPStatus(httpErr.Code)
			requestErr = errEchoHTTPError.SetInternalMessage(httpErr.Internal.Error())

		} else if _, ok := requestErr.(kitty.Reply); !ok {
			requestErr = errUnknownError.SetInternalMessage(requestErr.Error())
		}

		kerr := requestErr.(kitty.Reply)

		// Maybe error occur before set kitty context in middleware
		if kittyCtx, ok := c.Get(ContextKeyKittyCtx).(kitty.Context); ok {
			l = kittyCtx.Logger()
			t = kittyCtx.Translator()
		}

		msg, err := t.Translate(kerr.Key(), gutil.MapToKeyValue(kerr.Params())...)

		if err != nil {
			d := kerr.Data()
			d["__translation_err__"] = err.Error()
			kerr = kerr.SetData(d)

			msg = ""
		}

		// Report if need to report:
		if kerr.ShouldReport() {
			kerr.Report(l, t)
		}

		err = c.JSON(kerr.HTTPStatus(), kitty.NewBody(kerr.Code(), msg, kitty.Data(kerr.Data())))

		if err != nil {
			l.Error(err)
		}
	}

}
