package rabbitmq

import (
	"fmt"
	"net/http"

	"github.com/pegasus-cloud/msgqueue_client/msgqueue/utility"
)

// Update ...
func (q *QueueMethod) Update(qName, preDLQName, dlqName string) (err error) {
	_, predlq := parseMsgQueueName(preDLQName)
	if predlq != "" && preDLQName != dlqName {
		if err = q.unbindQueue(qName, preDLQName, qName+".dlx"); err != nil {
			return
		}
	}
	_, dlq := parseMsgQueueName(dlqName)
	if dlq != "" {
		if err = q.bindQueue(qName, dlqName, qName+".dlx"); err != nil {
			return
		}
	}

	return
}

func (q *QueueMethod) unbindQueue(qName, preDLQName, rkey string) (err error) {
	b, _, s, err := utility.SendRequest(
		"DELETE",
		q.getURL(fmt.Sprintf("bindings/%s/e/%s/q/%s/%s", q.Provider.AMQP.Vhost, q.Provider.AMQP.DeadLetterExchange, preDLQName, rkey)),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("[RabbitMQ](ename: %s, qname: %s, rkey: %s): %v", q.Provider.AMQP.DeadLetterExchange, preDLQName, qName, err)
	}

	if s != http.StatusNoContent {
		return fmt.Errorf("[RabbitMQ](ename: %s, qname: %s, rkey: %s): status: %d, message: %s", q.Provider.AMQP.DeadLetterExchange, preDLQName, qName, s, string(b))
	}
	return
}

func (q *QueueMethod) bindQueue(qName, dlqName, rkey string) (err error) {
	b, _, s, err := utility.SendRequest(
		"POST",
		q.getURL(fmt.Sprintf("bindings/%s/e/%s/q/%s", q.Provider.AMQP.Vhost, q.Provider.AMQP.DeadLetterExchange, dlqName)),
		map[string]string{
			ContentType: ApplicationJSON,
		},
		BindBody{
			RouteingKey: rkey,
		},
	)
	if err != nil {
		return fmt.Errorf("[RabbitMQ](ename: %s, qname: %s, rkey: %s): %v", q.Provider.AMQP.DeadLetterExchange, dlqName, qName, err)
	}

	if s != http.StatusCreated {
		return fmt.Errorf("[RabbitMQ](ename: %s, qname: %s, rkey: %s): status: %d, message: %s", q.Provider.AMQP.DeadLetterExchange, dlqName, qName, s, string(b))
	}

	return
}
