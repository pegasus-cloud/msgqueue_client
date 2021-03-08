package msgqueue

import (
	"fmt"

	"github.com/pegasus-cloud/msgqueue_client/msgqueue/common"
	"github.com/pegasus-cloud/msgqueue_client/msgqueue/rabbitmq"
)

// QueueFunc interface
type QueueFunc interface {
	common.QueueInterface
	common.MessageInterface
	common.RPCInterface
}

var provider Provider

// Provider ...
type Provider interface {
	Prepare()
	Connect() error
	Close()
	Queue() interface{}
	// GetChannel() (con, ch int, cha *amqp.Channel)
	// ReleaseChannel(con, ch int)
}

// Queue output queue method
func Queue() QueueFunc {
	return provider.Queue().(QueueFunc)
}

const (
	// RabbitMQEnum define enumeration of rabbitmq
	RabbitMQEnum = "rabbitmq"

	providerNameDoesNotSupport = "The ProviderName does not support currently, please select the other one"
)

// New make new message provider
func New(providerName string) {
	switch providerName {
	case RabbitMQEnum:
		provider = &rabbitmq.Provider{}
	default:
		panic(providerNameDoesNotSupport)
	}
	provider.Prepare()
	if err := provider.Connect(); err != nil {
		panic(err)
	}
}

// ManualNew input provider and connect to message queue service
func ManualNew(p Provider) {
	provider = p
	provider.Prepare()
	if err := provider.Connect(); err != nil {
		panic(err)
	}
}

// // GetChannel get channel in connection pool
// func GetChannel() (con, ch int, cha *amqp.Channel) {

// 	return provider.GetChannel()
// }

// // ReleaseChannel release channel back to connection pool
// func ReleaseChannel(con, ch int) {
// 	provider.ReleaseChannel(con, ch)
// 	return
// }

// Setup ...
func Setup() {
	if err := provider.Queue().(QueueFunc).Setup(); err != nil {
		fmt.Println(err)
	}
}

// Close ...
func Close() {
	provider.Close()
}
