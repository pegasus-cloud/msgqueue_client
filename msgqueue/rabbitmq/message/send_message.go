package message

import (
	"github.com/pegasus-cloud/msgqueue_client/rabbitmq/core"
)

// SendMessage send request to send message to the exchange
func (cfg *Config) SendMessage(name, rKey, mid, payload string, headers map[string]interface{}) (err error) {
	if err := core.GetAMQP().Publish(name, rKey, mid, payload, headers); err != nil {
		return err
	}

	return nil
}
