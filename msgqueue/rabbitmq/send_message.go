package rabbitmq

// SendMessage send request to send message to the exchange
func (q *QueueMethod) SendMessage(name, rKey, mid, payload string, headers map[string]interface{}) (err error) {
	if err := q.Provider.AMQP.Publish(name, rKey, mid, payload, headers); err != nil {
		return err
	}

	return nil
}
