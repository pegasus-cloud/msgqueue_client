package rabbitmq

import (
	"github.com/pegasus-cloud/msgqueue_client/msgqueue/common"
)

// ReceiveRPC receive rpc message
func (q *QueueMethod) ReceiveRPC(name string, del common.Delivering) (err error) {
	ch := q.Provider.AMQP.GetChannel()
	defer q.Provider.AMQP.ReleaseChannel(ch)

	msgs, err := ch.Consume(name, 1)
	if err != nil {
		return err
	}
	for d := range msgs {
		if err := q.Provider.AMQP.PublishRPC("", d.ReplyTo, "", d.CorrelationId, string(del(d.CorrelationId, d.Headers, d.Body))); err != nil {
			return err
		}
		d.Ack(false)
	}
	return nil
}
