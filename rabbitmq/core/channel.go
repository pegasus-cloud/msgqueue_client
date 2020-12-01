package core

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

var (
	chL sync.Mutex
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
func GetChannel() (con, ch int, cha *amqp.Channel) {
	chL.Lock()
	rand.Seed(time.Now().UnixNano())
	for {
		con = rand.Intn(globalConn.ConnectionNum)
		if len(idleChannel[con]) != 0 {
			break
		}
	}

	ch = globalConn.idleChannel[con][0]
	cha = globalConn.amqpConn[con].channel[ch]
	if len(idleChannel[con]) > 1 {
		globalConn.idleChannel[con] = globalConn.idleChannel[con][1:]
	} else {
		globalConn.idleChannel[con] = []int{}
	}
	chL.Unlock()
	return
}

// ReleaseChannel release channel resource
func ReleaseChannel(con, ch int) (err error) {
	chL.Lock()
	err = globalConn.amqpConn[con].channel[ch].Close()
	globalConn.idleChannel[con] = append(globalConn.idleChannel[con], ch)
	chL.Unlock()

	return
}

// GetAMQP return globalAMQP struct
func GetAMQP() *AMQP {
	return globalConn
}
