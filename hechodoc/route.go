package hechodoc

import (
	"fmt"
	"github.com/kamva/gutil"
	"github.com/labstack/echo/v4"
	"path"
	"regexp"
	"strings"
)

type RouteNameConverter interface {
	Tags(name string) []string // give tags from route name
	CamelCase(name string) string
}

// DefaultRouteNameConverter is the default route name Converter.
var DefaultRouteNameConverter = NewDividerNameConverter("::", 0)

// dividedNameConverter convert names which their format is just
// multiple words which joined with a div like , or :.
type dividedNameConverter struct {
	div string
	// specify maximum tags count which we can return. -1 means unlimited.
	maxTag int
}

func NewDividerNameConverter(div string, maxTag int) RouteNameConverter {
	return &dividedNameConverter{
		div:    div,
		maxTag: maxTag,
	}
}

func (c *dividedNameConverter) Tags(name string) []string {
	tags := strings.Split(name, c.div)
	if c.maxTag == -1 {
		return tags
	}
	return tags[:c.maxTag]
}

func (c *dividedNameConverter) CamelCase(name string) string {
	return camelCaseFromStringList(strings.Split(name, c.div))
}

var _ RouteNameConverter = &dividedNameConverter{}

// echoRawRouteNames returns string list of echo route names.
func echoRawRouteNames(e *echo.Echo) []string {
	routes := make([]string, len(e.Routes()))
	for i, r := range e.Routes() {
		routes[i] = r.Name
	}

	return routes
}

// OpenApiPathFromEchoPath convert echo path to openapi path.
// e.g, a/:id/:code => a/{id}/{code}
func OpenApiPathFromEchoPath(path string) string {
	for {
		colonIndex := strings.Index(path, ":")
		if colonIndex == -1 {
			return path
		}
		path = gutil.ReplaceRune(path, '{', colonIndex)

		slashPath := path[colonIndex:]
		slashIndex := strings.Index(slashPath, "/")
		if slashIndex == -1 {
			slashIndex = len(slashPath)
		}
		slashPath = gutil.ReplaceAt(slashPath, "}", slashIndex, slashIndex)
		path = path[:colonIndex] + slashPath
	}
}

type Route struct {
	BeginRouteVal string
	EndRouteVal   string
	Name          string
	RawName       string
	Method        string
	Path          string
	PathParams    []PathParam
	TagsString    string
	ParamsId      string
	SuccessRespId string
}

type PathParam struct {
	Name         string
	ExportedName string
}

func newRoute(r *echo.Route, c RouteNameConverter) Route {
	p := path.Join("/", OpenApiPathFromEchoPath(r.Path))
	return Route{
		BeginRouteVal: beginRouteVal(r.Name),
		EndRouteVal:   endRouteVal(r.Name),
		Name:          c.CamelCase(r.Name),
		RawName:       r.Name,
		Method:        r.Method,
		Path:          p,
		PathParams:    PathParams(p),
		TagsString:    strings.Join(c.Tags(r.Name), " "),
		ParamsId:      fmt.Sprintf("%sParams", c.CamelCase(r.Name)),
		SuccessRespId: fmt.Sprintf("%sSuccessResponse", c.CamelCase(r.Name)),
	}
}

func beginRouteVal(rName string) string {
	return BeginRoutePrefix + rName
}

func endRouteVal(rName string) string {
	return EndRoutePrefix + rName
}

var pathRegex = regexp.MustCompile("{.*?}")

func PathParams(p string) []PathParam {
	pList := pathRegex.FindAllString(p, -1)
	l := make([]PathParam, len(pList))
	for i, v := range pList {
		v = strings.Trim(v, "{")
		v = strings.Trim(v, "}")

		l[i] = PathParam{
			Name:         v,
			ExportedName: strings.ToUpper(v[:1]) + v[1:],
		}
	}
	return l
}
