package rabbitmq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pegasus-cloud/msgqueue_client/msgqueue/common"
	"github.com/pegasus-cloud/msgqueue_client/msgqueue/utility"
)

// Get ...
func (q *QueueMethod) Get(name string) (*common.Queue, error) {

	b, _, s, err := utility.SendRequest(
		"GET",
		q.getURL(fmt.Sprintf("queues/%s/%s", q.Provider.AMQP.Vhost, url.PathEscape(name))),
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("[SendRequest]: %v", err)
	}

	if s != http.StatusOK {
		return nil, fmt.Errorf("[RabbitMQ]: status: %d, message: %s", s, string(b))
	}

	queueInfo := &Queue{}
	if err := json.Unmarshal(b, queueInfo); err != nil {
		return nil, err
	}

	return parseQueueStruct(queueInfo), nil
}
