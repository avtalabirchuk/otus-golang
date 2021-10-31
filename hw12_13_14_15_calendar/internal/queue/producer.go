package queue

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

type Producer struct {
	connector    *Connector
	queueName    string
	exchangeType string
	doneCh       chan error
}

func NewProducer(uri, queueName, exchangeType string, maxReconnectAttempts, reconnectTimeoutMs int) *Producer {
	return &Producer{
		connector:    NewConnector(uri, queueName, exchangeType, maxReconnectAttempts, reconnectTimeoutMs),
		queueName:    queueName,
		exchangeType: exchangeType,
		doneCh:       make(chan error),
	}
}

func (p *Producer) GetDoneCh() <-chan error {
	return p.doneCh
}

func (p *Producer) Connect() error {
	return p.connector.Connect()
}

func (p *Producer) Publish(msg []byte) error {
	select {
	case <-p.connector.errCh:
		if errors.Is(p.connector.Reconnect(), ErrMaxConnectionAttempts) {
			p.doneCh <- ErrMaxConnectionAttempts
			// Not sure, but probably it makes sense to close it at all, otherwise, if any other process reads this error, it'll continue working
			// close(p.doneCh)
		}
	default:
		err := p.connector.GetChannel().Publish(
			"",          // exchange
			p.queueName, // routing key
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         msg,
			})
		if err != nil {
			return fmt.Errorf("error in Publishing: %w", err)
		}
	}
	return nil
}
