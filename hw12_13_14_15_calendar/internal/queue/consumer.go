package queue

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Consumer struct {
	connector        *Connector
	queueName        string
	exchangeType     string
	exchangeName     string
	bindingKey       string
	consumerTag      string
	qosPrefetchCount int
	done             chan error
}

func NewConsumer(uri, queueName, exchangeName, exchangeType, bindingKey, consumerTag string, qosPrefetchCount int, done chan error) *Consumer {
	return &Consumer{
		connector:        NewConnector(uri, exchangeName, exchangeType, done),
		queueName:        queueName,
		exchangeName:     exchangeName,
		exchangeType:     exchangeType,
		bindingKey:       bindingKey,
		consumerTag:      consumerTag,
		qosPrefetchCount: qosPrefetchCount,
		done:             done,
	}
}

func (c *Consumer) announceQueue() (<-chan amqp.Delivery, error) {
	channel := c.connector.GetChannel()
	queue, err := channel.QueueDeclare(
		c.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	err = channel.Qos(c.qosPrefetchCount, 0, false)
	if err != nil {
		return nil, fmt.Errorf("Error setting qos: %s", err)
	}

	if err = channel.QueueBind(
		queue.Name,
		c.bindingKey,
		c.exchangeName,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	msgs, err := channel.Consume(
		queue.Name,
		c.consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	return msgs, nil
}

func (c *Consumer) Handle(fn func(<-chan amqp.Delivery)) error {
	var err error
	if err = c.connector.Connect(); err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	msgs, err := c.announceQueue()
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	for {
		go fn(msgs)

		if <-c.done != nil {
			// msgs, err = c.reConnect()
			// if err != nil {
			// 	return fmt.Errorf("Reconnecting Error: %s", err)
			// }
		}
		fmt.Println("Reconnected... possibly")
	}
}
