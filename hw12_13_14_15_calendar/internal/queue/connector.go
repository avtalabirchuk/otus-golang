package queue

import (
	"errors"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Connector struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	uri          string
	exchangeType string
	exchangeName string
	done         chan error
}

func NewConnector(uri, exchangeName, exchangeType string, done chan error) *Connector {
	return &Connector{
		uri:          uri,
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		done:         done,
	}
}

func (c *Connector) GetChannel() *amqp.Channel {
	return c.channel
}

func (c *Connector) Connect() error {
	var err error

	if c.conn, err = amqp.Dial(c.uri); err != nil {
		return fmt.Errorf("dial: %s", err)
	}

	if c.channel, err = c.conn.Channel(); err != nil {
		return fmt.Errorf("channel: %s", err)
	}

	go func() {
		log.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
		c.done <- errors.New("channel closed")
	}()

	if err = c.channel.ExchangeDeclare(
		c.exchangeName,
		c.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("exchange declare: %s", err)
	}

	return nil
}
