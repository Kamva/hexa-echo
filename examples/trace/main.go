package main

import (
	hecho "github.com/kamva/hexa-echo"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Debug = true
	e.Logger = hecho.HexaToEchoLogger(hlog.NewPrinterDriver(), "debug")
	e.Use(hecho.Recover())
	e.HTTPErrorHandler = hecho.HTTPErrorHandler(hlog.NewPrinterDriver(), hexatranslator.NewEmptyDriver(), true)
	e.GET("/hi", func(c echo.Context) error {
		var a map[string]interface{}
		a["a"] = "12"
		return nil
	})
	e.Logger.Fatal(e.Start(":4444"))
}
