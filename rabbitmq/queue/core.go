package queue

import (
	"fmt"
	"strings"

	"github.com/pegasus-cloud/msgqueue_client/rabbitmq/core"
)

// Config define rabbitmq connection config
type Config struct {
	DeadLetterExchange string
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

func parseMsgQueueName(name string) (accountID, queueName string) {
	msgq := strings.Split(name, "@")

	return msgq[0], msgq[1]
}
