package hecho

import (
	"fmt"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"net/http"
)

// HTTPErrorHandler is the echo error handler.
// this function need to the HexaContext middleware.
func HTTPErrorHandler(l hexa.Logger, t hexa.Translator, debug bool) echo.HTTPErrorHandler {
	return func(rErr error, c echo.Context) {
		l := l
		t := t
		hlog.Debug("I'm handling error ---------------------")
		// We finally need to have a Reply or Error that internal error is stacked.
		stacked, baseErr := rErr, gutil.CauseErr(rErr)

		if httpErr, ok := baseErr.(*echo.HTTPError); ok {
			baseErr = errEchoHTTPError.SetHTTPStatus(httpErr.Code)
			hlog.Debug("is httpError")
			if httpErr.Code == http.StatusNotFound {
				hlog.Debug("is status not found")
				baseErr = errHTTPNotFoundError
				// NOTE: Do not set the "Internal" field of the http.StatusNotFound error.
				// otherwise for next 404 requests Echo checks if its internal error field
				// is not empty, it pass the internal field to this function as error instead
				// of real 404 error !

				httpErr.Message = fmt.Sprintf("route %s %s not found", c.Request().Method, c.Request().URL)
			}

			baseErr = baseErr.(hexa.Error).SetError(tracer.MoveStackIfNeeded(stacked, httpErr))
		} else {
			hlog.Debug("is not http error",hlog.String("base_error",fmt.Sprint(baseErr)))
			_, ok := baseErr.(hexa.Reply)
			_, ok2 := baseErr.(hexa.Error)

			if !ok && !ok2 {
				hlog.Debug("is not hexa error or reply")
				baseErr = errUnknownError.SetError(stacked)
			}
		}

		// Maybe error occur before set hexa context in middleware
		if hexaCtx, ok := c.Get(ContextKeyHexaCtx).(hexa.Context); ok {
			l = hexaCtx.Logger()
			t = hexaCtx.Translator()
		}

		if hexaErr, ok := baseErr.(hexa.Error); ok {
			handleError(hexaErr, c, l, t, debug)
		} else {
			handleReply(baseErr.(hexa.Reply), c, l, t)
		}
	}

}

func handleError(hexaErr hexa.Error, c echo.Context, l hexa.Logger, t hexa.Translator, debug bool) {
	msg, err := hexaErr.Localize(t)

	if err != nil {
		l.With(hlog.String("translation_key", hexaErr.ID())).Warn("translation for error id not found.")

		d := hexaErr.ReportData()
		d["_translation_err"] = err.Error()
		hexaErr = hexaErr.SetReportData(d)
	}

	// Report
	hexaErr.ReportIfNeeded(l, t)

	debugData := hexaErr.ReportData()
	debugData["err"] = hexaErr.Error()

	body := hexa.NewBody(hexaErr.ID(), msg, hexaErr.Data())

	body = body.Debug(debug, debugData)

	err = c.JSON(hexaErr.HTTPStatus(), body)

	if err != nil {
		l.Error("occurred error on request", hlog.Err(err))
	}
}

func handleReply(rep hexa.Reply, c echo.Context, l hexa.Logger, t hexa.Translator) {
	msg, err := t.Translate(rep.ID(), gutil.MapToKeyValue(rep.Data())...)

	if err != nil {
		l.With(hlog.String("translation_key", rep.ID())).Warn("translation for reply id not found.")
	}

	body := hexa.NewBody(rep.ID(), msg, rep.Data())

	err = c.JSON(rep.HTTPStatus(), body)

	if err != nil {
		l.Error("occurred error on request", hlog.Err(err))
	}
}
