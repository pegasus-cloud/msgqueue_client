package rabbitmq

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// SendRPC send rpc message and wait for receiving correlation id
func (q *QueueMethod) SendRPC(name, rKey, replyTo, payload string) (res string, err error) {
	ch := q.Provider.AMQP.GetChannel()
	defer q.Provider.AMQP.ReleaseChannel(ch)

	qname, msgs, err := ch.ExclusiveConsume(replyTo)
	if err != nil {
		return "", err
	}
	uid, _ := uuid.NewRandom()
	cid := uid.String()

	if err := q.Provider.AMQP.PublishRPC(name, rKey, qname, cid, payload); err != nil {
		return "", err
	}

	timer := time.NewTimer(time.Duration(q.Provider.AMQP.RPCTimeout) * time.Second)
	defer timer.Stop()

	for {
		select {
		case data := <-msgs:
			if data.CorrelationId == cid {
				res = string(data.Body)
				data.Ack(false)
				return res, nil
			}
		case <-timer.C:
			return "", errors.New("amqp rpc timeout")
		}
	}
}
