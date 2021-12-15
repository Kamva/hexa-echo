package hecho

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func emptyHandler(ctx echo.Context) error { return nil }

func TestExtractToken(t *testing.T) {
	sessionStore := sessions.NewCookieStore([]byte("secret-key"))
	sessionName := "token_test"
	tokenSessionKey := "token"
	cfg := ExtractTokenConfig{
		Skipper:                 middleware.DefaultSkipper,
		TokenContextKey:         AuthTokenContextKey,
		TokenLocationContextKey: AuthTokenLocationContextKey,
		Extractors: []TokenExtractor{
			CookieTokenExtractor(TokenCookieFieldAuthToken),
			HeaderTokenExtractor(TokenHeaderAuthorization),
			SessionTokenExtractor(sessionStore, sessionName, tokenSessionKey),
		},
	}
	testCases := []struct {
		Tag           string
		Cfg           ExtractTokenConfig
		Token         string
		TokenLocation TokenLocation
		RealTokenVal  string
	}{
		{Tag: "t1", Cfg: cfg, Token: "abc", TokenLocation: TokenLocationHeader, RealTokenVal: "bearer abc"},
		{Tag: "t2", Cfg: cfg, Token: "abc", TokenLocation: TokenLocationHeader, RealTokenVal: "bearer abc"},
		{Tag: "t3", Cfg: cfg, Token: "abc", TokenLocation: TokenLocationHeader, RealTokenVal: "token abc"},
		{Tag: "t4", Cfg: cfg, Token: "abc", TokenLocation: TokenLocationHeader, RealTokenVal: "TOKEN abc"},
		{Tag: "t5", Cfg: cfg, Token: "", TokenLocation: TokenLocationHeader, RealTokenVal: "Hi token"},
		{Tag: "t6", Cfg: cfg, Token: "", TokenLocation: TokenLocationHeader, RealTokenVal: "Bearer2 abc"},
		{Tag: "t7", Cfg: cfg, Token: "abc", TokenLocation: TokenLocationCookie, RealTokenVal: "abc"},
		{Tag: "t8", Cfg: cfg, Token: "", TokenLocation: TokenLocationCookie, RealTokenVal: ""},
		{Tag: "t9", Cfg: cfg, Token: "abc", TokenLocation: TokenLocationSession, RealTokenVal: "abc"},
		{Tag: "t10", Cfg: cfg, Token: "", TokenLocation: TokenLocationSession, RealTokenVal: ""},
	}

	e := echo.New()
	for _, item := range testCases {
		t.Run(item.Tag, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
			if item.TokenLocation == TokenLocationHeader && item.Token != "" {
				req.Header.Set(TokenHeaderAuthorization, item.RealTokenVal)
			} else if item.TokenLocation == TokenLocationCookie && item.Token != "" { // cookie
				req.Header.Set("Cookie", fmt.Sprintf("%s=%s;", TokenCookieFieldAuthToken, item.RealTokenVal))
			} else if item.Token != "" { // session
				sess, err := sessionStore.Get(req, sessionName)
				if assert.Nil(t, err) {
					sess.Values[tokenSessionKey] = item.RealTokenVal
				}
			}

			c := e.NewContext(req, httptest.NewRecorder())

			h := ExtractTokenWithConfig(item.Cfg)
			assert.Nil(t, h(emptyHandler)(c))

			if item.Token != "" {
				assert.Equal(t, item.Token, c.Get(item.Cfg.TokenContextKey))
				assert.Equal(t, item.TokenLocation, c.Get(item.Cfg.TokenLocationContextKey))
			} else {
				assert.Equal(t, nil, c.Get(item.Cfg.TokenContextKey))
				assert.Equal(t, nil, c.Get(item.Cfg.TokenLocationContextKey))
			}
		})
	}
}
