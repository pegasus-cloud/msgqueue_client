package queue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pegasus-cloud/ceph_client/ceph/utility"
	"github.com/pegasus-cloud/msgqueue_client/rabbitmq/core"
)

// GetQueue ...
func (cfg *Config) GetQueue(name string) (*Queue, error) {

	b, _, s, err := utility.SendRequest(
		"GET",
		cfg.getURL(fmt.Sprintf("queues/%s/%s", core.GetAMQP().Vhost, url.PathEscape(name))),
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

	return queueInfo, nil
}
