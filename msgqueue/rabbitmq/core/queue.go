package core

import "github.com/streadway/amqp"

// CreateExclusiveQueue ...
func (a *AMQP) CreateExclusiveQueue(rkey string, cha *amqp.Channel) (qname string, err error) {
	q, err := cha.QueueDeclare(
		rkey,  // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return "", err
	}

	return q.Name, nil
}
