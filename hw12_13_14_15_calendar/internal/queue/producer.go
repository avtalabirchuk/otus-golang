package queue

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Producer struct {
	connector    *Connector
	queueName    string
	exchangeName string
	exchangeType string
	routingKey   string
	done         chan error
}

func NewProducer(uri, queueName, routingKey, exchangeName, exchangeType string, done chan error) *Producer {
	return &Producer{
		connector:    NewConnector(uri, exchangeName, exchangeType, done),
		queueName:    queueName,
		routingKey:   routingKey,
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		done:         done,
	}
}

func (p *Producer) Handle(msgCh <-chan []byte, errCh chan<- error) error {
	var err error
	if err = p.connector.Connect(); err != nil {
		return fmt.Errorf("connection error: %v", err)
	}
	channel := p.connector.GetChannel()
	_, err = channel.QueueDeclare(
		p.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue declare: %s", err)
	}

	go func() {
		for {
			select {
			case msg := <-msgCh:
				err := channel.Publish(
					p.exchangeName, // exchange
					p.routingKey,   // routing key
					false,          // mandatory
					false,
					amqp.Publishing{
						DeliveryMode: amqp.Persistent,
						ContentType:  "application/json",
						Body:         msg,
					})
				if err != nil {
					errCh <- err
				}
			case <-p.done:
				// TODO: reconnect
				return
			}
		}
	}()
	return nil
}
