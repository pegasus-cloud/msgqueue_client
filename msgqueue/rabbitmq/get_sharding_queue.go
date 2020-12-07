package rabbitmq

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pegasus-cloud/msgqueue_client/msgqueue/common"
	"github.com/pegasus-cloud/msgqueue_client/msgqueue/utility"
)

// GetSharding ...
func (q *QueueMethod) GetSharding() (queuesInfo []*common.Queue, err error) {

	bs, err := q.bindingSource()
	if err != nil {
		return
	}
	queuesInfoMap, err := q.GetAll()
	if err != nil {
		return
	}
	for _, q := range bs {
		queue := &common.Queue{
			Name: q.Destination,
		}
		if val, ok := queuesInfoMap[q.Destination]; ok {
			queue.MessageNum = val.MessageNum
			queue.MessageSize = val.MessageSize
			queue.Memory = val.Memory
			queue.Consumers = val.Consumers
			queue.Node = val.Node
			queue.State = val.State
			queue.SlaveNodes = val.SlaveNodes
			queue.Policy = val.Policy
			queue.ReplicaSize = val.ReplicaSize
		}

		queuesInfo = append(queuesInfo, queue)
	}

	return
}
func (q *QueueMethod) bindingSource() ([]*BindingSources, error) {
	b, _, s, err := utility.SendRequest(
		"GET",
		q.getURL(fmt.Sprintf("exchanges/%s/%s/bindings/source", q.Provider.AMQP.Vhost, q.Provider.AMQP.ShardingQueueExchange)),
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("[RabbitMQ](ename: %s): %v", q.Provider.AMQP.ShardingQueueExchange, err)
	}

	if s != http.StatusOK {
		return nil, fmt.Errorf("[RabbitMQ](ename: %s): status: %d, message: %s", q.Provider.AMQP.ShardingQueueExchange, s, string(b))
	}
	bs := []*BindingSources{}
	if err := json.Unmarshal(b, &bs); err != nil {
		return nil, err
	}

	return bs, nil
}
