package queue

import (
	"fmt"
	"net/http"

	"github.com/pegasus-cloud/ceph_client/ceph/utility"
	"github.com/pegasus-cloud/msgqueue_client/rabbitmq/core"
)

// CreateQueue ...
func (cfg *Config) CreateQueue(name, dlqName string, len, ttl int) (err error) {

	if err = cfg.createQueue(name, len, ttl); err != nil {
		return
	}
	_, qName := parseMsgQueueName(dlqName)
	if qName == "" {
		return
	}
	err = cfg.bindExchange(cfg.Queue.DeadLetterExchange, dlqName, cfg.Queue.DeadLetterRoutingKey)

	return
}

func (cfg *Config) createQueue(name string, len, ttl int) error {

	b, _, s, err := utility.SendRequest(
		"PUT",
		cfg.getURL(fmt.Sprintf("queues/%s/%s", cfg.AMQP.Vhost, name)),
		map[string]string{
			core.ContentType: core.ApplicationJSON,
		},
		core.QueueBody{
			AutoDelete: false,
			Durable:    true,
			Arguments: core.QueueAgmt{
				MaxLength:            len,
				MessageTTL:           ttl * 1000,
				DeadLetterExchange:   cfg.Queue.DeadLetterExchange,
				DeadLetterRoutingKey: cfg.Queue.DeadLetterRoutingKey,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("[SendRequest](name: %s, len: %d, ttl: %d): %v", name, len, ttl, err)
	}

	if s != http.StatusCreated {
		return fmt.Errorf("[RabbitMQ](name: %s, len: %d, ttl: %d): status: %d, message: %s", name, len, ttl, s, string(b))
	}

	return nil
}

// Bind bind queue to the exchange
func (cfg *Config) bindExchange(ename, qname, rkey string) error {

	b, _, s, err := utility.SendRequest(
		"POST",
		cfg.getURL(fmt.Sprintf("bindings/%s/e/%s/q/%s", cfg.AMQP.Vhost, ename, qname)),
		map[string]string{
			core.ContentType: core.ApplicationJSON,
		},
		core.BindBody{
			RouteingKey: rkey,
		},
	)
	if err != nil {
		return fmt.Errorf("[RabbitMQ](ename: %s, qname: %s, rkey: %s): %v", ename, qname, rkey, err)
	}

	if s != http.StatusCreated {
		return fmt.Errorf("[RabbitMQ](ename: %s, qname: %s, rkey: %s): status: %d, message: %s", ename, qname, rkey, s, string(b))
	}

	return nil
}
