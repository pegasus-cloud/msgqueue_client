package rabbitmq

import (
	"errors"
	"time"
)

// SendRPC send rpc message and wait for receiving correlation id
func (q *QueueMethod) SendRPC(name, rKey, replyTo, payload string) (res string, err error) {

	t := time.Now()
	cid := "cid"
	qname, err := q.Provider.AMQP.CreateExclusiveQueue(replyTo)
	if err != nil {
		return "", err
	}
	if err := q.Provider.AMQP.PublishRPC(name, rKey, qname, cid, payload); err != nil {
		return "", err
	}

	msgs, err := q.Provider.AMQP.Consume(qname, 1)

	for {
		d := <-msgs
		if d.CorrelationId == cid {
			res = string(d.Body)
			return res, nil
		} else if int(time.Now().Sub(t).Seconds()) > q.Provider.AMQP.RPCTimeout {
			return "", errors.New("amqp rpc timeout")
		}
	}
}
