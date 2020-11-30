package core

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

func (ac *amqpConn) makeChannel(ch int) (err error) {
	cha, err := ac.connection.Channel()
	if err != nil {
		return fmt.Errorf("[RabbitMQ](): Failed to open a channel: %v", err)
	}

	if ac.chanNotify[ch] != nil {
		for range ac.chanNotify[ch] {
		}
	}
	ac.channel[ch] = cha
	ac.chanNotify[ch] = cha.NotifyClose(make(chan *amqp.Error))
	return nil
}

// GetChannel return channel of msgqueue
func GetChannel() *amqp.Channel {
	rand.Seed(time.Now().UnixNano())
	con := rand.Intn(globalConn.ConnectionNum)
	ch := rand.Intn(globalConn.ChannelNum)

	return globalConn.amqpConn[con].channel[ch]
}

// GetAMQP return globalAMQP struct
func GetAMQP() *AMQP {
	return globalConn
}
