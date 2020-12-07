package rabbitmq

import (
	"fmt"
	"net/http"

	"github.com/pegasus-cloud/msgqueue_client/msgqueue/utility"
)

// Create ...
func (q *QueueMethod) Create(name, dlqName string, len, ttl int) (err error) {

	if err = q.create(name, len, ttl); err != nil {
		return
	}
	_, qName := parseMsgQueueName(dlqName)
	if qName == "" {
		return
	}
	err = q.bindQueue(q.Provider.AMQP.DeadLetterExchange, dlqName, name+".dlx")

	return
}

func (q *QueueMethod) create(name string, len, ttl int) error {

	b, _, s, err := utility.SendRequest(
		"PUT",
		q.getURL(fmt.Sprintf("queues/%s/%s", q.Provider.AMQP.Vhost, name)),
		map[string]string{
			ContentType: ApplicationJSON,
		},
		QueueBody{
			AutoDelete: false,
			Durable:    true,
			Arguments: QueueAgmt{
				MaxLength:            len,
				MessageTTL:           ttl * 1000,
				DeadLetterExchange:   q.Provider.AMQP.DeadLetterExchange,
				DeadLetterRoutingKey: name + ".dlx",
			},
		},
	)
	if err != nil {
		return fmt.Errorf("[SendRequest](name: %s, len: %d, ttl: %d): %v", name, len, ttl, err)
	}

	if s != http.StatusCreated {
		return fmt.Errorf("[RabbitMQ](name: %s, len: %d, ttl: %d): status: %d, message: %s", name, len, ttl, s, string(b))
	}

	return nil
}
