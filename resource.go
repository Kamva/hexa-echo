package kecho

import "github.com/labstack/echo/v4"

type GetResource interface {
	Get(ctx echo.Context) error
}

type CreateResource interface {
	Create(ctx echo.Context) error
}
type UpdateResource interface {
	Update(ctx echo.Context) error
}

type PatchResource interface {
	Patch(ctx echo.Context) error
}

type DeleteResource interface {
	Delete(ctx echo.Context) error
}

// Resource define each method that exists in provided resource.
func Resource(group *echo.Group, resource interface{}) {
	if r, ok := resource.(GetResource); ok {
		group.GET("/:id", r.Get)
	}

	if r, ok := resource.(CreateResource); ok {
		group.POST("/:id", r.Create)
	}

	if r, ok := resource.(UpdateResource); ok {
		group.PUT("/:id", r.Update)
	}

	if r, ok := resource.(PatchResource); ok {
		group.PATCH("/:id", r.Patch)
	}

	if r, ok := resource.(DeleteResource); ok {
		group.DELETE("/:id", r.Delete)
	}
}
