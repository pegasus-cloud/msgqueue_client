package core

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

func makeChannel(con, ch int) (err error) {
	conn := globalConn.amqpConn[con]
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
	globalConn.idleChannel[con] = append(globalConn.idleChannel[con], ch)
	return nil
}

// GetChannel return channel of msgqueue
func GetChannel() (con, ch int, cha *amqp.Channel) {
	rand.Seed(time.Now().UnixNano())
	for {
		con = rand.Intn(globalConn.ConnectionNum)
		l.Lock()
		if len(globalConn.idleChannel[con]) != 0 {
			break
		}
		l.Unlock()
	}

	ch = globalConn.idleChannel[con][0]
	cha = globalConn.amqpConn[con].channel[ch]
	if len(globalConn.idleChannel[con]) > 1 {
		globalConn.idleChannel[con] = globalConn.idleChannel[con][1:]
	} else {
		globalConn.idleChannel[con] = []int{}
	}
	l.Unlock()
	return
}

// ReleaseChannel release channel resource
func ReleaseChannel(con, ch int) (err error) {
	err = globalConn.amqpConn[con].channel[ch].Close()

	return
}

// GetAMQP return globalAMQP struct
func GetAMQP() *AMQP {
	return globalConn
}
