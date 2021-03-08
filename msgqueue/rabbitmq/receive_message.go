package rabbitmq

import (
	"github.com/pegasus-cloud/msgqueue_client/msgqueue/common"
	"github.com/streadway/amqp"
)

// ReceiveMessage use amqp get and nack message to requeue messageuse. Checking function check deleted message and Receiving function receive message and append to output struct
func (q *QueueMethod) ReceiveMessage(name string, tgtsize int, del common.Checking, rec common.Receiving) (err error) {
	var (
		d    amqp.Delivery
		ok   bool
		size = 0
	)

	ch := q.Provider.AMQP.GetChannel()
	defer q.Provider.AMQP.ReleaseChannel(ch)
	_, _, cha := ch.GetInfo()

	for {
		d, ok, err = cha.Get(name, false)
		if !ok || size == tgtsize {
			d.Nack(true, true)
			break
		}
		if deleted, mesgAttr := del(d.MessageId); deleted {
			d.Ack(false)
		} else {
			d.Headers["Timestamp"] = d.Timestamp
			rec(d.MessageId, d.Headers, d.Body, mesgAttr)
			size++
			// d.Nack(false, true)
		}
	}

	return err
}
