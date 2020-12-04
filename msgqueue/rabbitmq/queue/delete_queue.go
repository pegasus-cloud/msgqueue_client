package queue

import (
	"fmt"
	"net/http"

	"github.com/pegasus-cloud/ceph_client/ceph/utility"
	"github.com/pegasus-cloud/msgqueue_client/rabbitmq/core"
)

// DeleteQueue ...
func (cfg *Config) DeleteQueue(name, dlqName string) (err error) {

	if err = cfg.deleteQueue(name); err != nil {
		return
	}
	_, qName := parseMsgQueueName(dlqName)
	if qName == "" {
		return
	}

	err = cfg.unBindExchange(cfg.DeadLetterExchange, dlqName, name+".dlx")
	if err != nil {
		return
	}

	return
}

// Delete Queue ...
func (cfg *Config) deleteQueue(name string) error {

	b, _, s, err := utility.SendRequest(
		"DELETE",
		cfg.getURL(fmt.Sprintf("queues/%s/%s", core.GetAMQP().Vhost, name)),
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

// UnBind unbind queue from the exchange
func (cfg *Config) unBindExchange(ename, qname, rkey string) error {

	b, _, s, err := utility.SendRequest(
		"DELETE",
		cfg.getURL(fmt.Sprintf("bindings/%s/e/%s/q/%s/%s", core.GetAMQP().Vhost, ename, qname, rkey)),
		map[string]string{
			core.ContentType: core.ApplicationJSON,
		},
		nil,
	)
	if err != nil {
		return fmt.Errorf("[RabbitMQ](ename: %s, qname: %s, rkey: %s): %v", ename, qname, rkey, err)
	}

	if s != http.StatusNoContent {
		return fmt.Errorf("[RabbitMQ](ename: %s, qname: %s, rkey: %s): status: %d, message: %s", ename, qname, rkey, s, string(b))
	}

	return nil
}
