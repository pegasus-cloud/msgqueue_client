package core

import (
	"fmt"

	"github.com/streadway/amqp"
)

// Consume ...
func (c *Channel) Consume(name string, contentSize int) (msgs <-chan amqp.Delivery, err error) {
	err = c.cha.Qos(
		contentSize, // prefetch count
		0,           // prefetch size
		false,       // global
	)
	if err != nil {
		return msgs, fmt.Errorf("[RabbitMQ](name: %s, contentSize: %d): Failed to set QoS: %v", name, contentSize, err)
	}

	msgs, err = c.cha.Consume(
		name,  // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return msgs, fmt.Errorf("[RabbitMQ](name: %s, contentSize: %d): Failed to register a consumer: %v", name, contentSize, err)
	}

	return msgs, nil
}

// ExclusiveConsume ...
func (c *Channel) ExclusiveConsume(rkey string) (qname string, msgs <-chan amqp.Delivery, err error) {

	q, err := c.cha.QueueDeclare(
		rkey,  // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return "", msgs, err
	}

	err = c.cha.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return "", msgs, fmt.Errorf("[RabbitMQ](name: %s): Failed to set QoS: %v", q.Name, err)
	}

	msgs, err = c.cha.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return "", msgs, fmt.Errorf("[RabbitMQ](name: %s): Failed to register a consumer: %v", q.Name, err)
	}

	return q.Name, msgs, nil
}
