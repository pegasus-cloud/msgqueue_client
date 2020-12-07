package rabbitmq

import (
	"fmt"
	"net/http"

	"github.com/pegasus-cloud/msgqueue_client/msgqueue/utility"
)

// Delete ...
func (q *QueueMethod) Delete(name, dlqName string) (err error) {

	if err = q.delete(name); err != nil {
		return
	}
	_, qName := parseMsgQueueName(dlqName)
	if qName == "" {
		return
	}

	err = q.unbindQueue(q.Provider.AMQP.DeadLetterExchange, dlqName, name+".dlx")
	if err != nil {
		return
	}

	return
}

func (q *QueueMethod) delete(name string) error {

	b, _, s, err := utility.SendRequest(
		"DELETE",
		q.getURL(fmt.Sprintf("queues/%s/%s", q.Provider.AMQP.Vhost, name)),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("[SendRequest](name: %s): %v", name, err)
	}

	if s != http.StatusNoContent {
		return fmt.Errorf("[RabbitMQ](name: %s): status: %d, message: %s", name, s, string(b))
	}

	return nil
}
