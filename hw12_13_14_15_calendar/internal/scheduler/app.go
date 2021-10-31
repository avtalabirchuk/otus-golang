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
	r             repository.Stats
	p             *queue.Producer
	scanTimeoutMs int
}

func New(r repository.Stats, p *queue.Producer, scanTimeoutMs int) *App {
	return &App{r: r, p: p, scanTimeoutMs: scanTimeoutMs}
}

func processEvents(repo repository.Stats, msgCh chan<- []byte) (err error) {
	events, err := repo.GetCurrentEvents()
	log.Debug().Msgf("Receiving events from DB %+v", events)
	if err != nil {
		return fmt.Errorf("GetCurrentEvents: %w", err)
	}
	if len(events) == 0 {
		// Skipping sending empty set
		return nil
	}
	bts, err := json.Marshal(events)
	if err != nil {
		return
	}
	if err := repo.MarkEventsAsProcessing(&events); err != nil {
		return fmt.Errorf("MarkEventsAsProcessing: %w", err)
	}
	msgCh <- bts
	return nil
}

func (app *App) scheduleEvents(msgCh chan []byte, errCh chan<- error) error {
	// nolint:durationcheck
	ticker := time.NewTicker(time.Duration(app.scanTimeoutMs) * time.Millisecond)
	for {
		select {
		case err := <-app.p.GetDoneCh():
			return err
		case <-ticker.C:
			log.Debug().Msg("Processing Events")
			if err := processEvents(app.r, msgCh); err != nil {
				errCh <- err
			}
			if err := app.r.DeleteObsoleteEvents(); err != nil {
				errCh <- fmt.Errorf("DeleteObsoleteEvents: %w", err)
			}
		}
	}
}

func (app *App) Run() error {
	msgCh := make(chan []byte)
	errCh := make(chan error)
	if err := app.p.Connect(); err != nil {
		return err
	}
	go func() {
		for msg := range msgCh {
			log.Debug().Msgf("Publishing message %s\n", msg)
			if err := app.p.Publish(msg); err != nil {
				log.Error().Msgf("Sending Error: %s\n", err)
			}
		}
	}()
	go func() {
		for b := range errCh {
			log.Error().Msgf("ERROR: %s\n", b)
		}
	}()
	return app.scheduleEvents(msgCh, errCh)
}
