package hechodoc

import (
	"github.com/labstack/echo/v4"
	"strings"
)

type RouteNameConverter interface {
	Tags(name string) []string // give tags from route name
	CamelCase(name string) string
}

// DefaultRouteNameConverter is the default route name Converter.
var DefaultRouteNameConverter = NewDividerNameConverter("::")

// dividedNameConverter convert names which their format is just
// multiple words which joined with a div like , or :.
type dividedNameConverter struct {
	div string
}

func NewDividerNameConverter(div string) RouteNameConverter {
	return &dividedNameConverter{
		div: div,
	}
}

func (c *dividedNameConverter) Tags(name string) []string {
	return strings.Split(name, c.div)
}

func (c *dividedNameConverter) CamelCase(name string) string {
	return camelCaseFromStringList(strings.Split(name, c.div))
}

var _ RouteNameConverter = &dividedNameConverter{}

// echoRawRouteNames returns string list of echo route names.
func echoRawRouteNames(e *echo.Echo)[]string{
	routes:=make([]string,len(e.Routes()))
	for i, r :=range e.Routes(){
		routes[i]= r.Name
	}

	return routes
}