package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

// ConsumeWithFunc define consume function
func (q *QueueMethod) ConsumeWithFunc(id, qname string, msgsFunc func(amqp.Delivery)) error {
	_, _, cha := q.Provider.AMQP.GetChannel()
	msgs, err := q.Provider.AMQP.Consume(qname, 1, cha)
	if err != nil {
		return fmt.Errorf("[RabbitMQ]%s: %v", qname, err)
	}

	for d := range msgs {
		msgsFunc(d)
		d.Ack(false)
	}

	return fmt.Errorf("[RabbitMQ]%s: consumer is closed", qname)
}
