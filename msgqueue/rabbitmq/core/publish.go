package core

import (
	"github.com/streadway/amqp"
)

// Publish ...
func (a *AMQP) Publish(ename, qname, mid, payload string, headers map[string]interface{}) (err error) {
	ch := a.GetChannel()
	defer a.ReleaseChannel(ch)
	err = ch.cha.Publish(ename, qname, false, false, amqp.Publishing{
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
	ch := a.GetChannel()
	defer a.ReleaseChannel(ch)
	err = ch.cha.Publish(ename, qname, false, false, amqp.Publishing{
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
