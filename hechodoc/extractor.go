package hechodoc

import (
	"bufio"
	"fmt"
	"github.com/kamva/gutil"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"text/template"
)

const (
	BeginRoutePrefix = "// route:begin: "
	EndRoutePrefix   = "// route:end: "
)

var DefaultSingleRouteTemplatePath = path.Join(gutil.SourcePath(), "default_template.tpl")
var beginRouteRegex = regexp.MustCompile(fmt.Sprintf("%s(.+)", BeginRoutePrefix))

func Placeholder(name string, prop string) string {
	return fmt.Sprintf("{{.%s.%s}}", name, prop)
}

func beginRouteVal(rName string) string {
	return BeginRoutePrefix + rName
}

func endRouteVal(rName string) string {
	return EndRoutePrefix + rName
}

func oldExtractedRoutes(f []byte) []string {
	allMatches := beginRouteRegex.FindAllStringSubmatch(string(f), -1)
	routes := make([]string, len(allMatches))
	for i, v := range allMatches {
		routes[i] = v[1]
	}
	return routes
}

type Extractor struct {
	echo           *echo.Echo
	singleRouteTpl *template.Template
	dst            string // Destination path
	converter      RouteNameConverter
}

type ExtractorOptions struct {
	Echo                    *echo.Echo
	ExtractDestinationPath  string // Destination path of extract filePath
	SingleRouteTemplatePath string
	Converter               RouteNameConverter
}

func NewExtractor(o ExtractorOptions) *Extractor {
	fName := path.Base(o.SingleRouteTemplatePath)
	functions := map[string]interface{}{
		"hold": Placeholder,
	}
	return &Extractor{
		echo:           o.Echo,
		singleRouteTpl: template.Must(template.New(fName).Funcs(functions).ParseFiles(o.SingleRouteTemplatePath)),
		dst:            o.ExtractDestinationPath,
		converter:      o.Converter,
	}
}

func (e *Extractor) Extract() error {
	file, err := os.OpenFile(e.dst, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return tracer.Trace(err)
	}
	defer file.Close()

	fBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return tracer.Trace(err)
	}

	buf := bufio.NewWriter(file)
	oldRoutes := oldExtractedRoutes(fBytes)

	// TODO: check if echo has repetitive route name, return error, we must not have any repetitive routes.
	// TODO: on each route check if route name is valid (using name converter), if its not valid => log and ignore it.
	// append new routes
	for _, r := range e.echo.Routes() {
		if !gutil.Contains(oldRoutes, r.Name) {
			if err := e.addRoute(r, buf); err != nil {
				return tracer.Trace(err)
			}
		}
	}
	return tracer.Trace(buf.Flush())
}

type GenerateRouteParams struct {
	Name          string // route name must be camelCase
	RawName       string // RawRoute is user's provided raw name for the route
	BeginRouteVal string
	EndRouteVal   string
}

func (e *Extractor) addRoute(r *echo.Route, w io.Writer) error {
	val := GenerateRouteParams{
		Name:          e.converter.CamelCase(r.Name),
		RawName:       r.Name,
		BeginRouteVal: beginRouteVal(r.Name),
		EndRouteVal:   endRouteVal(r.Name),
	}
	return e.singleRouteTpl.Execute(w, val)
}
