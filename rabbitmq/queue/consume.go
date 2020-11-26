package queue

import (
	"fmt"

	"github.com/streadway/amqp"
)

// ConsumeWithFunc define consume function
func (cfg *Config) ConsumeWithFunc(id, qname string, msgsFunc func(amqp.Delivery)) error {
	msgs, err := cfg.AMQP.Consume(qname, 1)
	if err != nil {
		return fmt.Errorf("[RabbitMQ]%s: %v", qname, err)
	}

	for d := range msgs {
		msgsFunc(d)
		d.Ack(false)
	}

	return fmt.Errorf("[RabbitMQ]%s: consumer is closed", qname)
}

type messageBody struct {
	Properties      struct{} `json:"properties"`
	RoutingKey      string   `json:"routing_key"`
	Payload         string   `json:"payload"`
	PayloadEncoding string   `json:"payload_encoding"`
}
