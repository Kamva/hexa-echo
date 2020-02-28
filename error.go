package kecho

import (
	"fmt"
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

		// We finally need to have a Reply or Error that internal error is stacked.
		stacked, baseErr := rErr, tracer.Cause(rErr)

		if httpErr, ok := baseErr.(*echo.HTTPError); ok {
			baseErr = errEchoHTTPError.SetHTTPStatus(httpErr.Code)

			if httpErr.Internal != nil {
				baseErr = errEchoHTTPError.SetError(tracer.MoveStack(stacked, httpErr.Internal))
			}

		} else {
			_, ok := baseErr.(kitty.Reply)
			_, ok2 := baseErr.(kitty.Error)
			fmt.Println(ok, ok2, baseErr, stacked)

			if !ok && !ok2 {
				baseErr = errUnknownError.SetError(stacked)
			}
		}

		// Maybe error occur before set kitty context in middleware
		if kittyCtx, ok := c.Get(ContextKeyKittyCtx).(kitty.Context); ok {
			l = kittyCtx.Logger()
			t = kittyCtx.Translator()
		}

		if kittyErr, ok := baseErr.(kitty.Error); ok {
			handleError(kittyErr, c, l, t, debug)
		} else {
			handleReply(baseErr.(kitty.Reply), c, l, t)
		}
	}

}

func handleError(kittyErr kitty.Error, c echo.Context, l kitty.Logger, t kitty.Translator, debug bool) {
	msg, err := t.Translate(kittyErr.Key(), gutil.MapToKeyValue(kittyErr.Params())...)

	if err != nil {
		l.WithFields("key", kittyErr.Key()).Warn("translation for specified key not found.")

		d := kittyErr.ReportData()
		d["__translation_err__"] = err.Error()
		kittyErr = kittyErr.SetReportData(d)
	}

	// Report
	kittyErr.ReportIfNeeded(l, t)

	debugData := kittyErr.ReportData()
	debugData["err"] = kittyErr.Error()

	body := kitty.NewBody(kittyErr.Code(), msg, kittyErr.Data())

	body = body.Debug(debug, debugData)

	err = c.JSON(kittyErr.HTTPStatus(), body)

	if err != nil {
		l.Error(err)
	}
}

func handleReply(rep kitty.Reply, c echo.Context, l kitty.Logger, t kitty.Translator) {
	msg, err := t.Translate(rep.Key(), gutil.MapToKeyValue(rep.Params())...)

	if err != nil {
		l.WithFields("key", rep.Key()).Warn("translation for specified key not found.")
	}

	body := kitty.NewBody(rep.Code(), msg, rep.Data())

	err = c.JSON(rep.HTTPStatus(), body)

	if err != nil {
		l.Error(err)
	}
}
