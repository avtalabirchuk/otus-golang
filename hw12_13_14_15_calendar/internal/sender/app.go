package sender

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/queue"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/repository"
)

type App struct {
	r           repository.Stats
	c           *queue.Consumer
	scanTimeout int
}

func New(r repository.Stats, c *queue.Consumer, scanTimeout int) *App {
	return &App{c: c, r: r, scanTimeout: scanTimeout}
}

func (app *App) processEvents(msg []byte) error {
	var events []repository.Event
	err := json.Unmarshal(msg, &events)
	if err != nil {
		log.Error().Msgf("cannot decode incoming message: %s", err)
	} else {
		for _, event := range events {
			fmt.Printf("[SEND] Event %s starts at %s, ends at %s\n", event.Title, event.StartDate.Format(time.RFC3339), event.EndDate.Format(time.RFC3339))
		}
		return app.r.MarkEventsAsSent(&events)
	}
	return nil
}

func (app *App) Run() error {
	return app.c.Handle(app.processEvents)
}
