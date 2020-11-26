package core

import (
	"sync"
	"time"
)

var (
	l sync.Mutex
)

// ReConnect ...
func (a *AMQP) ReConnect(con, ch int) {
	for {
		select {
		case <-a.amqpConn[con].connNotify:
		case <-a.amqpConn[con].chanNotify[ch]:
		case <-a.amqpConn[con].quit:
		}
	quit:
		for {
			select {
			case <-a.amqpConn[con].quit:
				return
			case <-a.amqpConn[con].connNotify:
				if ch < 0 {
					a.isClosed(con)
					// sleep 5s reconnect
					time.Sleep(time.Second * 5)
					l.Lock()
					if err := a.makeConnection(con); err != nil {
						l.Unlock()
						continue
					}
					l.Unlock()
				}
				break quit
			case <-a.amqpConn[con].chanNotify[ch]:
				a.isClosed(con, ch)
				// sleep 5s reconnect
				time.Sleep(time.Second * 5)
				l.Lock()
				if err := a.amqpConn[con].makeChannel(ch); err != nil {
					l.Unlock()
					continue
				}
				l.Unlock()
				break quit
			}
		}
	}
}

// IsClosed ...
func (a *AMQP) isClosed(con int, ch ...int) {
	if !a.amqpConn[con].connection.IsClosed() {
		if len(ch) == 1 {
			if a.amqpConn[con].channel[ch[0]] != nil {
				a.amqpConn[con].channel[ch[0]].Cancel("", false)
				a.amqpConn[con].channel[ch[0]].Close()
			}
		} else {
			for _, ch := range a.amqpConn[con].channel {
				ch.Cancel("", false)
				ch.Close()
			}

			// close message delivery
			a.amqpConn[con].connection.Close()
		}
	}
}
