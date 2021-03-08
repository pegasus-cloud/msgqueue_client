package rabbitmq

import (
	"fmt"
	"net/http"

	"github.com/pegasus-cloud/msgqueue_client/msgqueue/rabbitmq/core"
	"github.com/pegasus-cloud/msgqueue_client/msgqueue/utility"
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
	r.checkVhost()
	if err = r.AMQP.Connect(); err != nil {
		return
	}
	return nil
}

// // GetChannel get channel in connection pool
// func (r *Provider) GetChannel() (con, ch int, cha *amqp.Channel) {
// 	return r.AMQP.GetChannel()
// }

// // ReleaseChannel release channel back to connection pool
// func (r *Provider) ReleaseChannel(con, ch int) {
// 	r.AMQP.ReleaseChannel(con, ch)
// 	return
// }

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

func (r *Provider) checkVhost() error {
	b, _, s, err := utility.SendRequest(
		"PUT",
		r.Method.Queue.getURL(fmt.Sprintf("vhosts/%s", r.AMQP.Vhost)),
		map[string]string{
			ContentType: ApplicationJSON,
		},
		nil,
	)
	if err != nil {
		return err
	}

	if s != http.StatusNoContent && s != http.StatusCreated {
		return fmt.Errorf(string(b))
	}

	return nil

}
