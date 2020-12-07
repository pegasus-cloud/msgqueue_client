package rabbitmq

import (
	"github.com/pegasus-cloud/msgqueue_client/msgqueue/rabbitmq/core"

	"github.com/streadway/amqp"
)

// Provider ...
type Provider struct {
	core.AMQP
	Method *Method
}

// Method ...
type Method struct {
	Queue *QueueMethod
}

//QueueMethod ...
type QueueMethod struct {
	Provider *Provider
}

// Queue ...
func (r *Provider) Queue() interface{} {
	return r.Method.Queue
}

// Prepare provider
func (r *Provider) Prepare() {
	r.Method = &Method{
		Queue: &QueueMethod{
			Provider: r,
		},
	}
}

// Connect rabbitmq server
func (r *Provider) Connect() (err error) {
	if err = r.AMQP.Connect(); err != nil {
		return
	}
	return nil
}

// GetChannel get channel in connection pool
func (r *Provider) GetChannel() (con, ch int, cha *amqp.Channel) {
	return r.AMQP.GetChannel()
}

// ReleaseChannel release channel back to connection pool
func (r *Provider) ReleaseChannel(con, ch int) {
	r.AMQP.ReleaseChannel(con, ch)
	return
}

// Close the rabbitmq connection
func (r *Provider) Close() {
	r.AMQP.Close()
	return
}

// Setup the rabbitmq env
func (r *Provider) Setup() {
	r.Method.Queue.Setup()
	return
}
