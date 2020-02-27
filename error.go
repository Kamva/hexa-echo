package kecho

import (
	"github.com/Kamva/gutil"
	"github.com/Kamva/kitty"
	"github.com/Kamva/tracer"
	"github.com/labstack/echo/v4"
)

// HTTPErrorHandler is the echo error handler.
// this function need to the KittyContext middleware.
func HTTPErrorHandler(l kitty.Logger, t kitty.Translator, debug bool) echo.HTTPErrorHandler {
	return func(rErr error, c echo.Context) {
		l := l
		t := t

		if httpErr, ok := tracer.Cause(rErr).(*echo.HTTPError); ok {
			newErr := errEchoHTTPError.SetHTTPStatus(httpErr.Code)

			if httpErr.Internal != nil {
				newErr = errEchoHTTPError.SetInternalMessage(httpErr.Internal.Error())
			}

			rErr = tracer.MoveStack(rErr, newErr)

		} else if _, ok := tracer.Cause(rErr).(kitty.Reply); !ok {
			rErr = tracer.MoveStack(rErr, errUnknownError.SetInternalMessage(rErr.Error()))
		}

		kerr := tracer.Cause(rErr).(kitty.Reply)

		// Maybe error occur before set kitty context in middleware
		if kittyCtx, ok := c.Get(ContextKeyKittyCtx).(kitty.Context); ok {
			l = kittyCtx.Logger()
			t = kittyCtx.Translator()
		}

		msg, err := t.Translate(kerr.Key(), gutil.MapToKeyValue(kerr.Params())...)

		if err != nil {
			d := kerr.ReportData()
			d["__translation_err__"] = err.Error()
			kerr = kerr.SetReportData(d)

			msg = ""
		}

		// Report
		kerr.ReportIfNeeded(l, t)

		err = writeResponse(c, msg, kerr, debug)

		if err != nil {
			l.Error(err)
		}
	}

}

func writeResponse(c echo.Context, msg string, err kitty.Reply, debug bool) error {
	body := kitty.NewBody(err.Code(), msg, kitty.Data(err.Data()))

	debugData := kitty.Data(err.ReportData())

	if _, ok := err.Type().(kitty.ReplyTypeError); ok {
		debugData["err"] = err.Error()
	}

	body = body.Debug(debug, debugData)

	return c.JSON(err.HTTPStatus(), body)
}
