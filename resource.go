package hecho

import (
	"fmt"
	"github.com/Kamva/hexa"
	"github.com/labstack/echo/v4"
)

type (
	QueryResource interface {
		Query(ctx echo.Context) error
	}

	GetResource interface {
		Get(ctx echo.Context) error
	}

	CreateResource interface {
		Create(ctx echo.Context) error
	}

	UpdateResource interface {
		Update(ctx echo.Context) error
	}

	PatchResource interface {
		Patch(ctx echo.Context) error
	}

	DeleteResource interface {
		Delete(ctx echo.Context) error
	}

	Resource struct {
	}
)

// Ctx method extract the hexa context from the echo context.
func (r Resource) Ctx(c echo.Context) hexa.Context {
	return c.Get(ContextKeyHexaCtx).(hexa.Context)
}

// Resource define each method that exists in provided resource.
func ResourceAPI(group *echo.Group, resource interface{}, prefix string, m ...echo.MiddlewareFunc) {
	if r, ok := resource.(QueryResource); ok {
		group.GET("", r.Query, m...).Name = routeName(prefix, "query")
	}

	if r, ok := resource.(GetResource); ok {
		group.GET("/:id", r.Get, m...).Name = routeName(prefix, "get")
	}

	if r, ok := resource.(CreateResource); ok {
		group.POST("", r.Create, m...).Name = routeName(prefix, "create")
	}

	if r, ok := resource.(UpdateResource); ok {
		group.PUT("/:id", r.Update, m...).Name = routeName(prefix, "put")
	}

	if r, ok := resource.(PatchResource); ok {
		group.PATCH("/:id", r.Patch, m...).Name = routeName(prefix, "patch")
	}

	if r, ok := resource.(DeleteResource); ok {
		group.DELETE("/:id", r.Delete, m...).Name = routeName(prefix, "delete")
	}
}

func routeName(prefix, name string) string {
	return fmt.Sprintf("%s::%s", prefix, name)
}
