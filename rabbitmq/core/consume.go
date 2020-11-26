package core

import (
	"fmt"

	"github.com/streadway/amqp"
)

// Consume ...
func (a *AMQP) Consume(name string, contentSize int) (msgs <-chan amqp.Delivery, err error) {

	err = a.getChannel().Qos(
		contentSize, // prefetch count
		0,           // prefetch size
		false,       // global
	)
	if err != nil {
		return msgs, fmt.Errorf("[RabbitMQ](name: %s, contentSize: %d): Failed to set QoS: %v", name, contentSize, err)
	}
	msgs, err = a.getChannel().Consume(
		name,  // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		true,  // no-wait
		nil,   // args
	)
	if err != nil {
		return msgs, fmt.Errorf("[RabbitMQ](name: %s, contentSize: %d): Failed to register a consumer: %v", name, contentSize, err)
	}

	return msgs, nil
}
