package message

import (
	"fmt"

	"github.com/pegasus-cloud/msgqueue_client/rabbitmq/core"
	"github.com/streadway/amqp"
)

// ReceiveMessage use amqp consume and nack message to requeue message
func (cfg *Config) ReceiveMessage(name string, tgtsize int, del Checking, rec Receiving) (err error) {
	var (
		d    amqp.Delivery
		ok   bool
		size = 0
	)
	// msg.Nack(false, true)

	// msgs, err := core.GetAMQP().Consume(name, tgtsize)
	// if err != nil {
	// 	return err
	// }
	con, ch, cha := core.GetChannel()
	defer core.ReleaseChannel(con, ch)

	for {
		d, ok, err = cha.Get(name, false)
		fmt.Println(ok, size, tgtsize, err)
		if !ok || size == tgtsize {
			d.Nack(true, true)
			break
		}
		if deleted, mesgAttr := del(d.MessageId); deleted {
			d.Ack(false)
		} else {
			rec(d, mesgAttr)
			size++
			// d.Nack(false, true)
		}
	}

	return err
}