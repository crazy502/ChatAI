package rabbitmq

var (
	RMQMessage  *RabbitMQ
	RMQConsumer *RabbitMQ
)

func InitRabbitMQ() error {
	publisher, err := NewWorkRabbitMQ("Message")
	if err != nil {
		return err
	}

	consumer, err := NewWorkRabbitMQ("Message")
	if err != nil {
		publisher.Destroy()
		closeSharedConn()
		return err
	}

	RMQMessage = publisher
	RMQConsumer = consumer

	go RMQConsumer.Consume(MQMessage)
	return nil
}

func DestroyRabbitMQ() {
	if RMQConsumer != nil {
		RMQConsumer.Destroy()
		RMQConsumer = nil
	}

	if RMQMessage != nil {
		RMQMessage.Destroy()
		RMQMessage = nil
	}

	closeSharedConn()
}
