package common

import "github.com/streadway/amqp"

type QueueInterface interface {
	Create(name, dlqName string, len, ttl int) (err error)
	Get(name string) (*Queue, error)
	GetAll() (map[string]*Queue, error)
	GetSharding() (queues []*Queue, err error)
	Update(qName, preDLQName, dlqName string) (err error)
	Purge(name string) error
	Delete(name, dlqName string) (err error)
	Setup() error
	ConsumeWithFunc(id, qname string, msgsFunc func(amqp.Delivery)) error
}

type (
	// Queue ...
	Queue struct {
		Name         string
		Consumers    int
		ConsumersNum int
		MessageNum   int
		MessageSize  int
		Memory       int
		State        string
		Node         string
		SlaveNodes   []string
		Policy       string
		ReplicaSize  int
	}
)
