package hecho

import (
	"context"

	"github.com/kamva/hexa/sr"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
)

// EchoService implements hexa service.
type EchoService struct {
	*echo.Echo
	Address string
}

func (s *EchoService) Run() error {
	return tracer.Trace(s.Start(s.Address))
}

func (s *EchoService) Shutdown(ctx context.Context) error {
	return tracer.Trace(s.Echo.Shutdown(ctx))
}

var _ sr.Runnable = &EchoService{}
var _ sr.Shutdownable = &EchoService{}
