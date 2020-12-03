package queue

import (
	"fmt"
	"strings"

	"github.com/pegasus-cloud/msgqueue_client/rabbitmq/core"
)

type (
	// Config define rabbitmq connection config
	Config struct {
		DeadLetterExchange string
	}

	// Queue ...
	Queue struct {
		Name                      string   `json:"name"`
		Consumers                 int      `json:"consumers"`
		ConsumersNum              int      `json:"number_of_consumers"`
		MessageNum                int      `json:"messages"`
		MessageSize               int      `json:"message_bytes"`
		Memory                    int      `json:"memory"`
		State                     string   `json:"state"`
		Node                      string   `json:"node"`
		SlaveNodes                []string `json:"slave_nodes"`
		Policy                    string   `json:"policy"`
		EffectivePolicyDefinition struct {
			HAParams int `json:"ha-params"`
		} `json:"effective_policy_definition"`
	}
)

func (cfg *Config) getURL(endpoint string) string {
	return fmt.Sprintf(
		"http://%s:%s@%s:%d/api/%s",
		core.GetAMQP().Account,
		core.GetAMQP().Password,
		core.GetAMQP().IP,
		core.GetAMQP().HTTPPort,
		endpoint,
	)
}

func parseMsgQueueName(name string) (accountID, queueName string) {
	msgq := strings.Split(name, "@")

	return msgq[0], msgq[1]
}

func queuesInfo2Map(queues []*Queue) (maps map[string]*Queue) {
	maps = map[string]*Queue{}
	for _, q := range queues {
		maps[q.Name] = q
	}
	return
}
