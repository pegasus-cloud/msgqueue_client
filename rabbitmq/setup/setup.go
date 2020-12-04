package setup

import (
	"fmt"
	"net/http"

	"github.com/pegasus-cloud/msgqueue_client/rabbitmq/core"
	"github.com/pegasus-cloud/msgqueue_client/rabbitmq/utility"
	"github.com/streadway/amqp"
)

// Setup triton massage queue service
func (cfg *Config) Setup() error {
	if err := cfg.declareExchange(cfg.DeadLetterExchange, "direct", nil); err != nil {
		return err
	}
	if err := cfg.declareExchange(cfg.ShardingQueueExchange, "x-modulus-hash", nil); err != nil {
		return err
	}
	if err := cfg.declareExchange(cfg.DelayMessageExchange, "x-delayed-message", amqp.Table{"x-delayed-type": "x-modulus-hash"}); err != nil {
		return err
	}
	if err := cfg.bindExchange(cfg.DelayMessageExchange, cfg.ShardingQueueExchange, ""); err != nil {
		return err
	}

	if err := cfg.createPolicy(cfg.RabbitMQHaMode, "", "queues", core.Definition{
		HaMode:     "exactly",
		HaParam:    2,
		HaSyncMode: "automatic",
	}); err != nil {
		return err
	}

	if err := cfg.createPolicy(cfg.RabbitMQShardMode, "^sns.sharding$", "exchanges", core.Definition{
		ShardsPerNode: 2,
	}); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) declareExchange(ename, etype string, arg amqp.Table) error {
	con, ch, cha := core.GetChannel()
	defer core.ReleaseChannel(con, ch)

	return cha.ExchangeDeclare(
		ename, // name
		etype, // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		arg,   // arguments
	)
}

func (cfg *Config) bindExchange(src, dest, rkey string) error {

	con, ch, cha := core.GetChannel()
	defer core.ReleaseChannel(con, ch)

	return cha.ExchangeBind(
		dest,  // destination
		rkey,  // routing key
		src,   // source
		false, // no-wait
		nil,   // arguments
	)
}

func (cfg *Config) createPolicy(name, pattern, apply string, Definition core.Definition) error {
	b, _, s, err := utility.SendRequest(
		"PUT",
		cfg.getURL(fmt.Sprintf("policies/%s/%s", core.GetAMQP().Vhost, name)),
		map[string]string{
			core.ContentType: core.ApplicationJSON,
		},
		core.Permission{
			Pattern:    pattern,
			Definition: Definition,
			ApplyTo:    apply,
		},
	)
	if err != nil {
		return err
	}

	if s != http.StatusNoContent && s != http.StatusCreated {
		return fmt.Errorf(string(b))
	}

	return nil
}
