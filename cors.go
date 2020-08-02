package hecho

import (
	"github.com/kamva/hexa"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

// CORSConfig returns te CORS config.
func CorsConfig(config hexa.Config) middleware.CORSConfig {
	origins := config.GetList("CORS_ALLOW_ORIGINS")
	methods := config.GetList("CORS_ALLOW_METHODS")
	headers := config.GetList("CORS_ALLOW_HEADERS")
	if len(origins) == 0 {
		origins = []string{"*"}
	}
	if len(methods) == 0 {
		methods = []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete}
	}

	return middleware.CORSConfig{
		AllowOrigins: origins,
		AllowMethods: methods,
		AllowHeaders: headers,
	}
}
