package sender

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/queue"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/repository"
)

var ErrDecodeIncomingMessage = errors.New("cannot decode incoming message")

type App struct {
	c           *queue.Consumer
	scanTimeout int
}

func New(c *queue.Consumer, scanTimeout int) *App {
	return &App{c: c, scanTimeout: scanTimeout}
}

func (app *App) Send(events *[]repository.Event) {
	for _, event := range *events {
		fmt.Printf("[SEND] Event %s starts at %s, ends at %s\n", event.Title, event.StartDate.Format(time.RFC3339), event.EndDate.Format(time.RFC3339))
	}
}

func (app *App) Run(doneCh <-chan error) error {
	fn := func(msgCh <-chan amqp.Delivery) {
		for {
			select {
			case <-doneCh:
				return
			case msg := <-msgCh:
				var events []repository.Event
				err := json.Unmarshal(msg.Body, &events)
				if err != nil {
					log.Error().Msgf("%s: %s", ErrDecodeIncomingMessage, err)
				} else {
					app.Send(&events)
				}
			}
		}
	}
	return app.c.Handle(fn)
}
