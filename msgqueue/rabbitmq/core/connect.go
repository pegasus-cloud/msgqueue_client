package core

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/streadway/amqp"
)

type amqpConn struct {
	connection *amqp.Connection
	channel    map[int]*amqp.Channel

	quit       chan bool
	connNotify chan *amqp.Error
	chanNotify map[int]chan *amqp.Error
}

// AMQP define structure of amqp connection info
type AMQP struct {
	Account       string
	Password      string
	IP            string
	Port          int
	HTTPPort      int
	Timeout       int
	Vhost         string
	ConnectionNum int
	ChannelNum    int

	DeadLetterExchange    string
	ShardingQueueExchange string
	DelayMessageExchange  string

	amqpConn    map[int]*amqpConn
	idleChannel map[int][]int
}

// Connect rabbitmq server
func (a *AMQP) Connect() (err error) {
	a.amqpConn = make(map[int]*amqpConn)
	a.idleChannel = make(map[int][]int)
	for con := 0; con < a.ConnectionNum; con++ {
		a.amqpConn[con] = new(amqpConn)
		if err = a.connect(con); err != nil {
			return err
		}
	}

	return nil
}

func (a *AMQP) connect(con int) (err error) {
	if err = a.makeConnection(con); err != nil {
		return
	}
	go a.reConnect(con, -1)
	a.amqpConn[con].channel = make(map[int]*amqp.Channel)
	a.amqpConn[con].chanNotify = make(map[int]chan *amqp.Error)
	for ch := 0; ch < a.ChannelNum; ch++ {
		if err = a.makeChannel(con, ch); err != nil {
			return
		}
		go a.reConnect(con, ch)
	}

	return
}

// Close the rabbitmq connection
func (a *AMQP) Close() {
	for _, con := range a.amqpConn {
		close(con.quit)
	}
}

func (a *AMQP) makeConnection(con int) (err error) {
	conn, err := amqp.DialConfig(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%d/",
			url.PathEscape(a.Account),
			url.PathEscape(a.Password),
			a.IP,
			a.Port,
		),
		amqp.Config{
			Vhost: a.Vhost,
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, time.Duration(a.Timeout)*time.Second)
			},
		},
	)
	if err != nil {
		return err
	}

	if a.amqpConn[con].connNotify != nil {
		for range a.amqpConn[con].connNotify {
		}
	}
	a.amqpConn[con].connection = conn
	a.amqpConn[con].quit = make(chan bool)
	a.amqpConn[con].connNotify = conn.NotifyClose(make(chan *amqp.Error))
	a.idleChannel[con] = []int{}
	return nil
}
