package core

// CreateExclusiveQueue ...
func (a *AMQP) CreateExclusiveQueue(rkey string) (qname string, err error) {
	con, ch, cha := a.GetChannel()
	defer a.ReleaseChannel(con, ch)
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
