package queue

import (
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

type Connector struct {
	conn                 *amqp.Connection
	channel              *amqp.Channel
	uri                  string
	queueName            string
	exchangeType         string
	errCh                chan *amqp.Error
	reconnectAttempts    int
	maxReconnectAttempts int
	reconnectTimeoutMs   int
}

var ErrMaxConnectionAttempts = errors.New("maximum connection has been reached")

func NewConnector(uri, queueName, exchangeType string, maxReconnectAttempts, reconnectTimeoutMs int) *Connector {
	return &Connector{
		uri:                  uri,
		queueName:            queueName,
		exchangeType:         exchangeType,
		maxReconnectAttempts: maxReconnectAttempts,
		reconnectTimeoutMs:   reconnectTimeoutMs,
	}
}

func (c *Connector) GetChannel() *amqp.Channel {
	return c.channel
}

func (c *Connector) IsMaxConnError(err error) bool {
	return errors.Is(err, ErrMaxConnectionAttempts)
}

func (c *Connector) Connect() error {
	var err error

	if c.conn, err = amqp.Dial(c.uri); err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	if c.channel, err = c.conn.Channel(); err != nil {
		return fmt.Errorf("channel: %w", err)
	}

	_, err = c.channel.QueueDeclare(
		c.queueName, // queue
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)

	if err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}

	c.errCh = c.conn.NotifyClose(make(chan *amqp.Error))

	go func() {
		log.Info().Msgf("closing: %s", <-c.errCh)
	}()

	return nil
}

func (c *Connector) Reconnect() error {
	c.reconnectAttempts++
	if c.reconnectAttempts > c.maxReconnectAttempts {
		return ErrMaxConnectionAttempts
	}
	log.Info().Msg("Connection to queue has been failed. Trying to reconnect...")
	// nolint:durationcheck
	time.Sleep(time.Duration(c.reconnectTimeoutMs) * time.Millisecond)
	if err := c.Connect(); err != nil {
		return err
	}
	c.reconnectAttempts = 0
	return nil
}
