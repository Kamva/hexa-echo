package kecho

import (
	"github.com/Kamva/elogrus/v4"
	"github.com/Kamva/kitty"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// KittyLoggerToEchoLogger convert kitty logger to echo logger.
func KittyLoggerToEchoLogger(logger kitty.Logger) echo.Logger {
	switch logger.Core().(type) {
	case *logrus.Entry:
		return elogrus.GetEchoLogger(logger.Core().(*logrus.Entry))
	}

	return nil
}
