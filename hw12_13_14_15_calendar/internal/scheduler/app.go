package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

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

func processEvents(repo repository.Stats, msgCh chan<- []byte) (err error) {
	events, err := repo.GetCurrentEvents()
	if err != nil {
		return fmt.Errorf("GetCurrentEvents: %v", err)
	}
	bts, err := json.Marshal(events)
	if err != nil {
		return
	}
	if err := repo.MarkEventsAsProcessing(&events); err != nil {
		return fmt.Errorf("MarkEventsAsProcessing: %v", err)
	}
	msgCh <- bts
	return nil
}

func (app *App) Run(doneCh <-chan error) error {
	msgCh := make(chan []byte)
	errCh := make(chan error)
	if err := app.p.Handle(msgCh, errCh); err != nil {
		return err
	}
	ticker := time.NewTicker(time.Duration(app.scanTimeout) * time.Second)
	go func() {
		for b := range errCh {
			log.Error().Msgf("ERROR: %s\n", b)
		}
	}()
	for {
		select {
		case <-doneCh:
			return nil
		case <-ticker.C:
			go func() {
				err := processEvents(app.r, msgCh)
				if err != nil {
					errCh <- err
				}
			}()
			go func() {
				err := app.r.DeleteObsoleteEvents()
				if err != nil {
					errCh <- fmt.Errorf("DeleteObsoleteEvents: %v", err)
				}
			}()
		}
	}
}
