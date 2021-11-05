package queue

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

type Consumer struct {
	connector        *Connector
	queueName        string
	exchangeType     string
	qosPrefetchCount int
}

func NewConsumer(uri, queueName, exchangeType string, qosPrefetchCount, maxReconnectAttempts, reconnectTimeoutMs int) *Consumer {
	return &Consumer{
		connector:        NewConnector(uri, queueName, exchangeType, maxReconnectAttempts, reconnectTimeoutMs),
		queueName:        queueName,
		exchangeType:     exchangeType,
		qosPrefetchCount: qosPrefetchCount,
	}
}

func (c *Consumer) Connect() error {
	return c.connector.Connect()
}

func (c *Consumer) IsMaxConnError(err error) bool {
	return c.connector.IsMaxConnError(err)
}

func (c *Consumer) Consume() (<-chan amqp.Delivery, error) {
	channel := c.connector.GetChannel()
	err := channel.Qos(c.qosPrefetchCount, 0, false)
	if err != nil {
		return nil, fmt.Errorf("error setting qos: %s", err)
	}
	return channel.Consume(
		c.queueName, // queue
		"",          // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
}

func (c *Consumer) Reconnect() (deliveries <-chan amqp.Delivery, err error) {
	for {
		cErr := c.connector.Reconnect()
		if cErr != nil && c.IsMaxConnError(cErr) {
			return nil, cErr
		}
		deliveries, err = c.Consume()
		if err == nil {
			break
		}
	}
	return
}

func (c *Consumer) Handle(fn func([]byte) error) error {
	if err := c.Connect(); err != nil {
		return err
	}
	deliveries, err := c.Consume()
	if err != nil {
		return err
	}
	for {
		select {
		case msg := <-deliveries:
			log.Debug().Msgf("Receiving events from queue: %+v", msg.Body)
			err := fn(msg.Body)
			if err != nil {
				log.Error().Msgf("Processing event: %s", err)
			}
		case <-c.connector.errCh:
			deliveries, err = c.Reconnect()
			if err != nil {
				return err
			}
		}
	}
}
