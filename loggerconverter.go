package hecho

import (
	"github.com/Kamva/elogrus/v4"
	"github.com/Kamva/hexa"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// HexaToEchoLogger convert hexa logger to echo logger.
func HexaToEchoLogger(logger hexa.Logger) echo.Logger {
	switch logger.Core().(type) {
	case *logrus.Entry:
		return elogrus.GetEchoLogger(logger.Core().(*logrus.Entry))
	}

	return nil
}
