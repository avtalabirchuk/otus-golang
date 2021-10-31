package queue

import (
	"fmt"
)

func GetRabbitMQURL(args ...interface{}) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", args...)
}
