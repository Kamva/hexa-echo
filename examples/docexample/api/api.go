package api

import (
	"github.com/kamva/gutil"
	"github.com/labstack/echo/v4"
	"net/http"
)

type HiRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type HiResponse struct {
	Code string       `json:"code"`
	Data ResponseData `json:"data"`
}

type ResponseData struct {
	Say string `json:"say"`
}

func RegisterRoutes(e *echo.Echo) {
	e.GET("/hi", func(c echo.Context) error {
		resp := HiResponse{
			Code: "hi.success_get",
			Data: ResponseData{
				Say: "hiiii",
			},
		}

		return c.JSON(http.StatusOK, gutil.StructToMap(resp))
	}).Name = "hi::say"

	e.POST("/hi/:id/:code", func(c echo.Context) error {
		resp := HiResponse{
			Code: "hi.success_create",
			Data: ResponseData{
				Say: "object created :)",
			},
		}

		return c.JSON(http.StatusOK, gutil.StructToMap(resp))
	}).Name = "hi::create"
}
