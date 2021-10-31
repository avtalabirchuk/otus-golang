package app

import (
	"errors"
	"sync"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/repository"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/service"
)

var ErrUnrecognizedServiceType = errors.New("cannot create service, because type was not recognized. Supported types: http, grpc")

type App struct {
	r repository.Base
	c *config.Config
}

func New(c *config.Config, r repository.Base) (*App, error) {
	return &App{c: c, r: r}, nil
}

func (app *App) Run(errCh *chan error, doneCh *chan bool) {
	s := service.New(app.r)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if app.c.GRPCAddress != "" {
			*errCh <- s.RunGRPC(app.c.GRPCAddress)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if app.c.HTTPAddress != "" && app.c.GRPCAddress != "" {
			*errCh <- s.RunHTTP(app.c.GRPCAddress, app.c.HTTPAddress)
		}
	}()
	go func() {
		wg.Wait()
		close(*doneCh)
	}()
}
