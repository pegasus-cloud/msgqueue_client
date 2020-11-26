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

func (a *AMQP) getChannel() *amqp.Channel {
	rand.Seed(time.Now().UnixNano())
	con := rand.Intn(a.ConnectionNum)
	ch := rand.Intn(a.ChannelNum)

	return a.amqpConn[con].channel[ch]
}
