package rabbitmq

import (
	"errors"
	"fmt"
	"time"
)

// SendRPC send rpc message and wait for receiving correlation id
func (q *QueueMethod) SendRPC(name, rKey, replyTo, payload string) (res string, err error) {
	con, ch, cha := q.Provider.AMQP.GetChannel()
	qname, err := q.Provider.AMQP.CreateExclusiveQueue(replyTo, cha)
	if err != nil {
		return "", err
	}
	cid := qname
	if err := q.Provider.AMQP.PublishRPC(name, rKey, qname, cid, payload); err != nil {
		return "", err
	}
	msgs, err := q.Provider.AMQP.Consume(qname, 1, cha)
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-msgs:
			if data.CorrelationId == cid {
				res = string(data.Body)
				data.Ack(false)
				q.Provider.AMQP.ReleaseChannel(con, ch)
				return res, nil
			}
		case <-time.After(time.Duration(q.Provider.AMQP.RPCTimeout) * time.Second):
			return "", errors.New("amqp rpc timeout")
		}
	}
}
