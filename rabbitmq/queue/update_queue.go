package queue

import (
	"fmt"
	"net/http"

	"github.com/pegasus-cloud/ceph_client/ceph/utility"
	"github.com/pegasus-cloud/msgqueue_client/rabbitmq/core"
)

// UpdateQueue ...
func (cfg *Config) UpdateQueue(qName, preDLQName, dlqName string) (err error) {
	_, predlq := parseMsgQueueName(preDLQName)
	if predlq != "" && preDLQName != dlqName {
		if err = cfg.unbindQueue(qName, preDLQName); err != nil {
			return
		}
	}
	_, dlq := parseMsgQueueName(dlqName)
	if dlq != "" {
		if err = cfg.bindQueue(qName, dlqName); err != nil {
			return
		}
	}

	return
}

func (cfg *Config) unbindQueue(qName, preDLQName string) (err error) {
	b, _, s, err := utility.SendRequest(
		"DELETE",
		cfg.getURL(fmt.Sprintf("bindings/%s/e/%s/q/%s/%s", core.GetAMQP().Vhost, cfg.DeadLetterExchange, preDLQName, qName+".dlx")),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("[RabbitMQ](ename: %s, qname: %s, rkey: %s): %v", cfg.DeadLetterExchange, preDLQName, qName, err)
	}

	if s != http.StatusNoContent {
		return fmt.Errorf("[RabbitMQ](ename: %s, qname: %s, rkey: %s): status: %d, message: %s", cfg.DeadLetterExchange, preDLQName, qName, s, string(b))
	}
	return
}

func (cfg *Config) bindQueue(qName, dlqName string) (err error) {
	b, _, s, err := utility.SendRequest(
		"POST",
		cfg.getURL(fmt.Sprintf("bindings/%s/e/%s/q/%s", core.GetAMQP().Vhost, cfg.DeadLetterExchange, dlqName)),
		map[string]string{
			core.ContentType: core.ApplicationJSON,
		},
		core.BindBody{
			RouteingKey: qName + ".dlx",
		},
	)
	if err != nil {
		return fmt.Errorf("[RabbitMQ](ename: %s, qname: %s, rkey: %s): %v", cfg.DeadLetterExchange, dlqName, qName, err)
	}

	if s != http.StatusCreated {
		return fmt.Errorf("[RabbitMQ](ename: %s, qname: %s, rkey: %s): status: %d, message: %s", cfg.DeadLetterExchange, dlqName, qName, s, string(b))
	}

	return
}
