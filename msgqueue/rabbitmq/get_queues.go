package rabbitmq

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pegasus-cloud/msgqueue_client/msgqueue/common"
	"github.com/pegasus-cloud/msgqueue_client/msgqueue/utility"
)

// GetAll ...
func (q *QueueMethod) GetAll() (map[string]*common.Queue, error) {
	b, _, s, err := utility.SendRequest(
		"GET",
		q.getURL(fmt.Sprintf("queues/")),
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("[SendRequest]: %v", err)
	}

	if s != http.StatusOK {
		return nil, fmt.Errorf("[RabbitMQ]: status: %d, message: %s", s, string(b))
	}

	amqpQueues := []*Queue{}
	if err := json.Unmarshal(b, &amqpQueues); err != nil {
		return nil, err
	}

	queuesInfo := []*common.Queue{}
	for _, q := range amqpQueues {
		queuesInfo = append(queuesInfo, parseQueueStruct(q))
	}

	return queuesInfo2Map(queuesInfo), nil
}
