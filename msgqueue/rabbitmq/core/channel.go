package core

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

func (a *AMQP) makeChannel(con, ch int) (err error) {
	conn := a.amqpConn[con]
	cha, err := conn.connection.Channel()
	if err != nil {
		return fmt.Errorf("[RabbitMQ](): Failed to open a channel: %v", err)
	}

	if conn.chanNotify[ch] != nil {
		for range conn.chanNotify[ch] {
		}
	}
	conn.channel[ch] = cha
	conn.chanNotify[ch] = cha.NotifyClose(make(chan *amqp.Error))
	a.idleChannel[con] = append(a.idleChannel[con], ch)
	return nil
}

// GetChannel return channel of msgqueue
func (a *AMQP) GetChannel() (con, ch int, cha *amqp.Channel) {
	rand.Seed(time.Now().UnixNano())
	for {
		con = rand.Intn(a.ConnectionNum)
		l.Lock()
		if len(a.idleChannel[con]) != 0 {
			break
		}
		l.Unlock()
	}

	ch = a.idleChannel[con][0]
	cha = a.amqpConn[con].channel[ch]
	if len(a.idleChannel[con]) > 1 {
		a.idleChannel[con] = a.idleChannel[con][1:]
	} else {
		a.idleChannel[con] = []int{}
	}
	l.Unlock()
	return
}

// ReleaseChannel release channel resource
func (a *AMQP) ReleaseChannel(con, ch int) (err error) {
	err = a.amqpConn[con].channel[ch].Close()

	return
}
