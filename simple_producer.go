package sarama

// SimpleProducer publishes Kafka messages. It routes messages to the correct broker, refreshing metadata as appropriate,
// and parses responses for errors. You must call Close() on a producer to avoid leaks, it may not be garbage-collected automatically when
// it passes out of scope (this is in addition to calling Close on the underlying client, which is still necessary).
type SimpleProducer struct {
	producer *Producer
	topic    string
}

// NewSimpleProducer creates a new SimpleProducer using the given client and topic.
func NewSimpleProducer(client *Client, topic string) (*SimpleProducer, error) {
	if topic == "" {
		return nil, ConfigurationError("Empty topic")
	}

	config := new(ProducerConfig)
	config.RequiredAcks = 1
	config.AckSuccesses = true

	prod, err := NewProducer(client, config)

	if err != nil {
		return nil, err
	}

	return &SimpleProducer{prod, topic}, nil
}

// SendMessage produces a message with the given key and value. The partition to send to is selected
// at randome. To send strings as either key or value, see the StringEncoder type.
func (sp *SimpleProducer) SendMessage(key, value Encoder) error {
	sp.producer.SendMessage(sp.topic, key, value)

	err := <-sp.producer.Errors() // we always get something because AckSuccesses is true

	if err != nil {
		return err.Err
	}

	return nil
}

// Close shuts down the producer and flushes any messages it may have buffered. You must call this function before
// a producer object passes out of scope, as it may otherwise leak memory. You must call this before calling Close
// on the underlying client.
func (sp *SimpleProducer) Close() error {
	return sp.producer.Close()
}
