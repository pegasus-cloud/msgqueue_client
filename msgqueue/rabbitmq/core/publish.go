package core

import (
	"github.com/streadway/amqp"
)

// Publish ...
func (a *AMQP) Publish(ename, qname, mid, payload string, headers map[string]interface{}) (err error) {
	con, ch, cha := a.GetChannel()
	defer a.ReleaseChannel(con, ch)
	err = cha.Publish(ename, qname, false, false, amqp.Publishing{
		Headers:      headers,
		MessageId:    mid,
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(payload),
	})
	if err != nil {
		return err
	}

	return nil
}

// PublishRPC ...
func (a *AMQP) PublishRPC(ename, qname, replyTo, cid, payload string) (err error) {
	con, ch, cha := a.GetChannel()
	defer a.ReleaseChannel(con, ch)
	err = cha.Publish(ename, qname, false, false, amqp.Publishing{
		ReplyTo:       replyTo,
		CorrelationId: cid,
		DeliveryMode:  amqp.Persistent,
		ContentType:   "text/plain",
		Body:          []byte(payload),
	})
	if err != nil {
		return err
	}

	return nil
}
