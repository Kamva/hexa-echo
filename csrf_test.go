package hecho

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCSRFSkipperByAuthTokenLocation(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/abc", nil)
	rec := httptest.NewRecorder()

	table := []struct {
		Tag           string
		TokenLocation TokenLocation
		Skip          bool
	}{
		{"t1", TokenLocationUnknown, false},
		{"t1", TokenLocationHeader, true},
		{"t1", TokenLocationCookie, false},
		{"t1", TokenLocationSession, false},
	}

	for _, item := range table {
		t.Run(item.Tag, func(t *testing.T) {
			c := e.NewContext(req, rec)
			c.Set(AuthTokenLocationContextKey,item.TokenLocation)
			assert.Equal(t, item.Skip,CSRFSkipperByAuthTokenLocation(c))
		})
	}
}
