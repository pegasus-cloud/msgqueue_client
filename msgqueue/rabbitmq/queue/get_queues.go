package queue

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pegasus-cloud/ceph_client/ceph/utility"
)

// GetQueues ...
func (cfg *Config) GetQueues() (map[string]*Queue, error) {
	b, _, s, err := utility.SendRequest(
		"GET",
		cfg.getURL(fmt.Sprintf("queues/")),
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("[SendRequest]: %v", err)
	}

	if s != http.StatusOK {
		return nil, fmt.Errorf("[RabbitMQ]: status: %d, message: %s", s, string(b))
	}

	queuesInfo := []*Queue{}
	if err := json.Unmarshal(b, &queuesInfo); err != nil {
		return nil, err
	}

	return queuesInfo2Map(queuesInfo), nil
}
