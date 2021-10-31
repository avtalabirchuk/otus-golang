package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/queue"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/repository"
)

var ErrUnrecognizedServiceType = errors.New("cannot create service, because type was not recognized. Supported types: http, grpc")

type App struct {
	r           repository.Stats
	p           *queue.Producer
	scanTimeout int
}

func New(r repository.Stats, p *queue.Producer, scanTimeout int) *App {
	return &App{r: r, p: p, scanTimeout: scanTimeout}
}

func (app *App) Run(doneCh <-chan error) error {
	msgCh := make(chan []byte)
	errCh := make(chan error)
	if err := app.p.Handle(msgCh, errCh); err != nil {
		return err
	}
	ticker := time.NewTicker(time.Duration(app.scanTimeout) * time.Second)
	go func() {
		for {
			select {
			case b := <-errCh:
				fmt.Printf("ERRRR: %s\n", b)
			}
		}
	}()
	for {
		select {
		case <-doneCh:
			return nil
		case <-ticker.C:
			events, err := app.r.GetCurrentEvents()
			if err != nil {
				errCh <- err
			} else {
				b, err := json.Marshal(events)
				if err != nil {
					errCh <- err
				} else {
					msgCh <- b
					// err := app.r.MarkEventsAsProcessing(&events)
					// if err != nil {
					// 	errCh <- err
					// }
				}
			}
		}
	}
}
