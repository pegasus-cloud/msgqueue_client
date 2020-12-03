package queue

import (
	"fmt"
	"net/http"

	"github.com/pegasus-cloud/ceph_client/ceph/utility"
	"github.com/pegasus-cloud/msgqueue_client/rabbitmq/core"
)

// PurgeQueue ...
func (cfg *Config) PurgeQueue(name string) error {

	b, _, s, err := utility.SendRequest(
		"DELETE",
		cfg.getURL(fmt.Sprintf("queues/%s/%s/contents", core.GetAMQP().Vhost, name)),
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
