package mq

var (
	RMQMessage  *RabbitMQ
	RMQConsumer *RabbitMQ
)

func InitRabbitMQ() error {
	publisher, err := NewWorkRabbitMQ(MessageQueueName)
	if err != nil {
		return err
	}

	consumer, err := NewWorkRabbitMQ(MessageQueueName)
	if err != nil {
		publisher.Destroy()
		closeSharedConn()
		return err
	}

	RMQMessage = publisher
	RMQConsumer = consumer
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
