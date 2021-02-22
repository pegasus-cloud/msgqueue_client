package rabbitmq

import (
	"errors"
	"time"

	"github.com/pegasus-cloud/msgqueue_client/msgqueue/common"
)

// ReceiveRPC receive rpc message
func (q *QueueMethod) ReceiveRPC(name string, rec common.Receiving) (err error) {

	msgs, err := q.Provider.AMQP.Consume(name, 1)
	select {
	case d := <-msgs:
		rec(d.CorrelationId, d.Headers, d.Body, nil)
		if err := q.Provider.AMQP.PublishRPC("", d.ReplyTo, "", d.CorrelationId, Success); err != nil {
			return err
		}
	case <-chTimeout(q.Provider.AMQP.RPCTimeout):
		return errors.New("amqp rpc timeout")
	default:
		return errors.New("unexpected error")
	}

	return nil
}

func chTimeout(timeout int) chan bool {
	tch := make(chan bool, 1)

	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		tch <- true
	}()

	return tch
}
