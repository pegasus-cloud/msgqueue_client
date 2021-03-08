package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

// ConsumeWithFunc define consume function
func (q *QueueMethod) ConsumeWithFunc(id, qname string, msgsFunc func(amqp.Delivery)) error {
	ch := q.Provider.AMQP.GetChannel()
	defer q.Provider.AMQP.ReleaseChannel(ch)
	msgs, err := ch.Consume(qname, 1)
	if err != nil {
		return fmt.Errorf("[RabbitMQ]%s: %v", qname, err)
	}

	for d := range msgs {
		msgsFunc(d)
		d.Ack(false)
	}

	return fmt.Errorf("[RabbitMQ]%s: consumer is closed", qname)
}
