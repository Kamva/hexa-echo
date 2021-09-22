package hecho

import (
	"errors"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
)

// Write writes reply as response.
// You MUST have hexa context in your echo context to use
// this function to use its logger and translator.
func Write(c echo.Context, reply hexa.Reply) error {
	hexaCtx := Ctx(c)
	if hexaCtx == nil {
		return tracer.Trace(errors.New("invalid hexa context, we can not write reply as response"))
	}
	l := hexaCtx.Logger()
	t := hexaCtx.Translator()

	return WriteWithOpts(c, l, t, reply)
}

// WriteWithOpts writes the reply as response.
func WriteWithOpts(c echo.Context, l hexa.Logger, t hexa.Translator, reply hexa.Reply) error {
	msg, err := t.Translate(reply.ID(), gutil.MapToKeyValue(reply.Data())...)
	if err != nil {
		l.With(hlog.String("translation_key", reply.ID())).Warn("translation for reply id not found.")
	}

	body := hexa.NewBody(reply.ID(), msg, reply.Data())
	err = c.JSON(reply.HTTPStatus(), body)
	if err != nil {
		l.Error("occurred error on request", hlog.Err(err))
	}
	return tracer.Trace(err)
}
