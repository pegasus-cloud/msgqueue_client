package setup

import (
	"fmt"

	"github.com/pegasus-cloud/msgqueue_client/rabbitmq/core"
)

// Config ...
type Config struct {
	DeadLetterExchange    string
	ShardingQueueExchange string
	DelayMessageExchange  string
	RabbitMQHaMode        string
	RabbitMQShardMode     string
}

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
