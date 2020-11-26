package queue

import (
	"fmt"
	"strings"

	"github.com/pegasus-cloud/msgqueue_client/rabbitmq/core"
)

// Config define rabbitmq connection config
type Config struct {
	core.AMQP
	Queue
}

// Queue ...
type Queue struct {
	DeadLetterExchange   string
	DeadLetterRoutingKey string
}

func (cfg *Config) getURL(endpoint string) string {
	return fmt.Sprintf(
		"http://%s:%s@%s:%d/api/%s",
		cfg.AMQP.Account,
		cfg.AMQP.Password,
		cfg.AMQP.IP,
		cfg.AMQP.HTTPPort,
		endpoint,
	)
}

func parseMsgQueueName(name string) (accountID, queueName string) {
	msgq := strings.Split(name, "@")

	return msgq[0], msgq[1]
}
