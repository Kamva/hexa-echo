package kecho

import (
	"github.com/Kamva/kitty"
	"github.com/labstack/echo/v4"
)

type (
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

// Ctx method extract the kitty context from the echo context.
func (r Resource) Ctx(c echo.Context) kitty.Context {
	return c.Get(ContextKeyKittyCtx).(kitty.Context)

}

// Resource define each method that exists in provided resource.
func ResourceAPI(group *echo.Group, resource interface{}, m ...echo.MiddlewareFunc) {
	if r, ok := resource.(GetResource); ok {
		group.GET("/:id", r.Get, m...)
	}

	if r, ok := resource.(CreateResource); ok {
		group.POST("/:id", r.Create, m...)
	}

	if r, ok := resource.(UpdateResource); ok {
		group.PUT("/:id", r.Update, m...)
	}

	if r, ok := resource.(PatchResource); ok {
		group.PATCH("/:id", r.Patch, m...)
	}

	if r, ok := resource.(DeleteResource); ok {
		group.DELETE("/:id", r.Delete, m...)
	}
}
