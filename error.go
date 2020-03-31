package hecho

import (
	"fmt"
	"github.com/Kamva/gutil"
	"github.com/Kamva/hexa"
	"github.com/Kamva/tracer"
	"github.com/labstack/echo/v4"
	"net/http"
)

// HTTPErrorHandler is the echo error handler.
// this function need to the HexaContext middleware.
func HTTPErrorHandler(l hexa.Logger, t hexa.Translator, debug bool) echo.HTTPErrorHandler {
	return func(rErr error, c echo.Context) {
		l := l
		t := t

		// We finally need to have a Reply or Error that internal error is stacked.
		stacked, baseErr := rErr, tracer.Cause(rErr)

		if httpErr, ok := baseErr.(*echo.HTTPError); ok {
			baseErr = errEchoHTTPError.SetHTTPStatus(httpErr.Code)
			if httpErr.Code == http.StatusNotFound {
				baseErr = errHTTPNotFoundError
				httpErr.Internal=fmt.Errorf("route %s %s not found",c.Request().Method,c.Request().URL)
			}
			if httpErr.Internal != nil {
				baseErr = baseErr.(hexa.Error).SetError(tracer.MoveStack(stacked, httpErr.Internal))
			}

		} else {
			_, ok := baseErr.(hexa.Reply)
			_, ok2 := baseErr.(hexa.Error)

			if !ok && !ok2 {
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
	msg, err := t.Translate(hexaErr.Key(), gutil.MapToKeyValue(hexaErr.Params())...)

	if err != nil {
		l.WithFields("key", hexaErr.Key()).Warn("translation for specified key not found.")

		d := hexaErr.ReportData()
		d["__translation_err__"] = err.Error()
		hexaErr = hexaErr.SetReportData(d)
	}

	// Report
	hexaErr.ReportIfNeeded(l, t)

	debugData := hexaErr.ReportData()
	debugData["err"] = hexaErr.Error()

	body := hexa.NewBody(hexaErr.Code(), msg, hexaErr.Data())

	body = body.Debug(debug, debugData)

	err = c.JSON(hexaErr.HTTPStatus(), body)

	if err != nil {
		l.Error(err)
	}
}

func handleReply(rep hexa.Reply, c echo.Context, l hexa.Logger, t hexa.Translator) {
	msg, err := t.Translate(rep.Key(), gutil.MapToKeyValue(rep.Params())...)

	if err != nil {
		l.WithFields("key", rep.Key()).Warn("translation for specified key not found.")
	}

	body := hexa.NewBody(rep.Code(), msg, rep.Data())

	err = c.JSON(rep.HTTPStatus(), body)

	if err != nil {
		l.Error(err)
	}
}
