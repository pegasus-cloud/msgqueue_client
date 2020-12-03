package queue

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pegasus-cloud/ceph_client/ceph/utility"
	"github.com/pegasus-cloud/msgqueue_client/rabbitmq/core"
	"go.uber.org/zap"
)

// GetSharding ...
func (cfg *Config) GetSharding() (queues []*Queue, err error) {

	bs, err := cfg.bindingSource()
	if err != nil {
		return
	}
	queuesInfoMap, err := cfg.GetQueues()
	if err != nil {
		return
	}

	for _, q := range bs {
		queue := &Queue{
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
			queue.EffectivePolicyDefinition.HAParams = val.EffectivePolicyDefinition.HAParams
		}

		queues = append(queues, queue)
	}

	return
}
func (cfg *Config) bindingSource() ([]*core.BindingSources, error) {
	b, _, s, err := utility.SendRequest(
		"GET",
		cfg.getURL(fmt.Sprintf("exchanges/%s/%s/bindings/source", core.GetAMQP().Vhost, cfg.ShardingQueueExchange)),
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("[RabbitMQ](ename: %s): %v", cfg.ShardingQueueExchange, err)
	}

	if s != http.StatusOK {
		return nil, fmt.Errorf("[RabbitMQ](ename: %s): status: %d, message: %s", cfg.ShardingQueueExchange, s, string(b))
	}
	bs := []*core.BindingSources{}
	if err := json.Unmarshal(b, &bs); err != nil {
		zap.L().With(zap.String("json", "Unmarshal"), zap.String("data", string(b))).Error(err.Error())
		return nil, err
	}

	return bs, nil
}
