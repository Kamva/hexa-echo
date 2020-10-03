package main

import (
	hecho "github.com/kamva/hexa-echo"
	"github.com/kamva/hexa/hlog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/url"
	"path"
)

const (
	dstUrl = "https://webhook.site"
	prefix = "hook"
)

func main() {
	e := echo.New()
	e.Group(prefix).Use(reverseProxy())
	e.Logger.Fatal(e.Start(":6006"))
	// e.g., go to localhost:6006/hook/{your path} to proxy it to the webhook.test/{your path}
}

func reverseProxy() echo.MiddlewareFunc {
	// Setup proxy
	url1, err := url.Parse(dstUrl)
	if err != nil {
		hlog.Error(err)
	}

	targets := []*hecho.ProxyTarget{{URL: url1}}
	proxyCfg := hecho.ProxyConfig{
		Skipper:              middleware.DefaultSkipper,
		Balancer:             hecho.NewRoundRobinBalancer(targets),
		ContextKey:           "target",
		Rewrite:              map[string]string{path.Join(prefix, "*"): "$1"}, // stack/* => $1
		ReverseProxyProvider: hecho.NewSingleHostReverseProxyRewriteHost,
	}

	return hecho.ProxyWithConfig(proxyCfg)
}
