package common

import "github.com/streadway/amqp"

type (
	MessageInterface interface {
		SendMessage(name, rKey, mid, payload string, headers map[string]interface{}) (err error)
		ReceiveMessage(name string, tgtsize int, del Checking, rec Receiving) (err error)
	}
	// Checking ...
	Checking func(messageIDs string) (deleted bool, iMetadata []byte)

	// Receiving ...
	Receiving func(d amqp.Delivery, iMetadata []byte)
)
