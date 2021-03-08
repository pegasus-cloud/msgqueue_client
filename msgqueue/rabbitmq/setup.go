package rabbitmq

import (
	"fmt"
	"net/http"

	"github.com/pegasus-cloud/msgqueue_client/msgqueue/utility"
	"github.com/streadway/amqp"
)

// Setup triton massage queue service
func (q *QueueMethod) Setup() error {
	if err := q.declareExchange(q.Provider.AMQP.DeadLetterExchange, "direct", nil); err != nil {
		return fmt.Errorf("declare exchange failed: name: %s, type: %s, message: %v", q.Provider.AMQP.DeadLetterExchange, "direct", err)
	}
	if err := q.declareExchange(q.Provider.AMQP.ShardingQueueExchange, "x-modulus-hash", nil); err != nil {
		return fmt.Errorf("declare exchange failed: name: %s, type: %s, message: %v", q.Provider.AMQP.ShardingQueueExchange, "x-modulus-hash", err)
	}
	if err := q.declareExchange(q.Provider.AMQP.DelayMessageExchange, "x-delayed-message", amqp.Table{"x-delayed-type": "x-modulus-hash"}); err != nil {
		return fmt.Errorf("declare exchange failed: name: %s, type: %s, message: %v", q.Provider.AMQP.DelayMessageExchange, "x-delayed-message", err)
	}
	if err := q.bindExchange(q.Provider.AMQP.DelayMessageExchange, q.Provider.AMQP.ShardingQueueExchange, ""); err != nil {
		return fmt.Errorf("bind exchange failed: name: %s to %s, message: %v", q.Provider.AMQP.DelayMessageExchange, q.Provider.AMQP.ShardingQueueExchange, err)
	}

	if err := q.createPolicy(RabbitMQHaMode, "", "queues", Definition{
		HaMode:     "exactly",
		HaParam:    2,
		HaSyncMode: "automatic",
	}); err != nil {
		return fmt.Errorf("create policy failed: mode: %s, message: %v", RabbitMQHaMode, err)
	}

	if err := q.createPolicy(RabbitMQShardMode, "^sns.sharding$", "exchanges", Definition{
		ShardsPerNode: 2,
	}); err != nil {
		return fmt.Errorf("create policy failed: mode: %s, message: %v", RabbitMQShardMode, err)
	}

	return nil
}

func (q *QueueMethod) declareExchange(ename, etype string, arg amqp.Table) error {
	ch := q.Provider.AMQP.GetChannel()
	defer q.Provider.AMQP.ReleaseChannel(ch)
	_, _, amqpChan := ch.GetInfo()

	return amqpChan.ExchangeDeclare(
		ename, // name
		etype, // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		arg,   // arguments
	)
}

func (q *QueueMethod) bindExchange(src, dest, rkey string) error {
	ch := q.Provider.AMQP.GetChannel()
	defer q.Provider.AMQP.ReleaseChannel(ch)
	_, _, amqpChan := ch.GetInfo()

	return amqpChan.ExchangeBind(
		dest,  // destination
		rkey,  // routing key
		src,   // source
		false, // no-wait
		nil,   // arguments
	)
}

func (q *QueueMethod) createPolicy(name, pattern, apply string, Definition Definition) error {
	b, _, s, err := utility.SendRequest(
		"PUT",
		q.getURL(fmt.Sprintf("policies/%s/%s", q.Provider.AMQP.Vhost, name)),
		map[string]string{
			ContentType: ApplicationJSON,
		},
		Permission{
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
