package main

import (
	"github.com/kamva/gutil"
	hecho "github.com/kamva/hexa-echo"
	"github.com/kamva/hexa-echo/examples/docexample/api"
	_ "github.com/kamva/hexa-echo/examples/docexample/doc/gen"
	"github.com/kamva/hexa-echo/hechodoc"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"github.com/labstack/echo/v4"
	"log"
	"os"
	"path"
)

func boot() *echo.Echo {
	e := echo.New()
	e.Debug = true
	e.Logger = hecho.HexaToEchoLogger(hlog.NewPrinterDriver(), "debug")
	e.Use(hecho.Recover())
	e.HTTPErrorHandler = hecho.HTTPErrorHandler(hlog.NewPrinterDriver(), hexatranslator.NewEmptyDriver(), true)
	api.RegisterRoutes(e)
	return e
}

var converter = hechodoc.DefaultRouteNameConverter
var extractPath = path.Join(gutil.SourcePath(), "doc/openapi_tpl.go")
var renderPath = path.Join(gutil.SourcePath(), "doc/gen/openapi_autogen.go")

func main() {
	// Echo instance
	if len(os.Args) < 2 {
		log.Fatal("provide action name please")
	}
	action := os.Args[1]
	if action == "" {
		log.Fatal("provide action please.")
	}
	switch action {
	case "server":
		runServer()
	case "extract":
		extract()
	case "render":
		renderTemplate()
	case "trim":
		trim()
	default:
		log.Fatal("unknown action")
	}
}

func runServer() {
	e := boot()
	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func extract() {
	extractor := hechodoc.NewExtractor(hechodoc.ExtractorOptions{
		Echo:                    boot(),
		ExtractDestinationPath:  extractPath,
		SingleRouteTemplatePath: hechodoc.DefaultSingleRouteTemplatePath,
		Converter:               converter,
	})
	gutil.PanicErr(extractor.Extract())
}

func trim() {
	trimmer := hechodoc.NewTrimmer(hechodoc.TrimmerOptions{
		Echo:                   boot(),
		ExtractDestinationPath: extractPath,
	})

	gutil.PanicErr(trimmer.Trim())
}

func renderTemplate() {
	renderer := hechodoc.NewRenderer(hechodoc.RendererOptions{
		Echo:              boot(),
		ExtractedFilePath: extractPath,
		Destination:       renderPath,
		Converter:         converter,
	})
	gutil.PanicErr(renderer.Render())
}
