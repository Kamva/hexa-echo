package hechodoc

import (
	"fmt"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"strings"
	"text/template"

	"os"
)

type Route struct {
	Name          string
	Method        string
	Path          string
	TagsString    string
	ParamsId      string
	SuccessRespId string
}

func newRoute(r *echo.Route, c RouteNameConverter) Route {

	return Route{
		Name:          c.CamelCase(r.Name),
		Method:        r.Method,
		Path:          r.Path,
		TagsString:    strings.Join(c.Tags(r.Name)," "),
		ParamsId:      fmt.Sprintf("%sParams", c.CamelCase(r.Name)),
		SuccessRespId: fmt.Sprintf("%sSuccessResponse", c.CamelCase(r.Name)),
	}
}


type Renderer struct {
	echo      *echo.Echo
	src       *template.Template // src path
	dst       string             // dst path
	converter RouteNameConverter
}

type RendererOptions struct {
	Echo              *echo.Echo
	ExtractedFilePath string
	Destination       string
	Converter         RouteNameConverter
}

func NewRenderer(o RendererOptions) *Renderer {
	return &Renderer{
		echo:      o.Echo,
		src:       template.Must(template.ParseFiles(o.ExtractedFilePath)),
		dst:       o.Destination,
		converter: o.Converter,
	}
}

func (r *Renderer) Render() error {
	file, err := os.OpenFile(r.dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return tracer.Trace(err)
	}
	defer file.Close()

	var routes = make(map[string]Route)

	for _, route := range r.echo.Routes() {
		routes[r.converter.CamelCase(route.Name)] = newRoute(route, r.converter)
	}
	return tracer.Trace(r.src.Execute(file, routes))
}
