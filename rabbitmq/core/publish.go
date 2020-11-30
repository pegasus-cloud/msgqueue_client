package core

import (
	"github.com/streadway/amqp"
)

// Publish ...
func (a *AMQP) Publish(ename, qname, mid, payload string, headers map[string]interface{}) (err error) {

	err = GetChannel().Publish(ename, qname, false, false, amqp.Publishing{
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
