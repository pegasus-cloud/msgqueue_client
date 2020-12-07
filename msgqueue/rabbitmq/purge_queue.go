package rabbitmq

import (
	"fmt"
	"net/http"

	"github.com/pegasus-cloud/msgqueue_client/msgqueue/utility"
)

// Purge ...
func (q *QueueMethod) Purge(name string) error {

	b, _, s, err := utility.SendRequest(
		"DELETE",
		q.getURL(fmt.Sprintf("queues/%s/%s/contents", q.Provider.AMQP.Vhost, name)),
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
