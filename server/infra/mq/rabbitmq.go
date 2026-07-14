package mq

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"server/infra/config"
	"server/pkg/observe"

	"github.com/streadway/amqp"
)

var (
	conn   *amqp.Connection
	connMu sync.Mutex
)

func initConn() error {
	connMu.Lock()
	defer connMu.Unlock()

	if conn != nil {
		return nil
	}

	cfg := config.GetConfig()
	mqURL := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		cfg.RabbitMQUsername,
		cfg.RabbitMQPassword,
		cfg.RabbitMQHost,
		cfg.RabbitMQPort,
		cfg.RabbitMQVhost,
	)

	connection, err := amqp.Dial(mqURL)
	if err != nil {
		return fmt.Errorf("rabbitmq connection failed: %w", err)
	}

	conn = connection
	return nil
}

func closeSharedConn() {
	connMu.Lock()
	defer connMu.Unlock()

	if conn != nil {
		_ = conn.Close()
		conn = nil
	}
}

type RabbitMQ struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	Exchange string
	Key      string
	mu       sync.Mutex
	confirms <-chan amqp.Confirmation
}

func NewRabbitMQ(exchange, key string) *RabbitMQ {
	return &RabbitMQ{Exchange: exchange, Key: key}
}

func (r *RabbitMQ) Destroy() {
	if r == nil {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.channel != nil {
		_ = r.channel.Close()
		r.channel = nil
	}
	r.conn = nil
}

func NewWorkRabbitMQ(queue string) (*RabbitMQ, error) {
	rabbitMQ := NewRabbitMQ("", queue)

	if err := initConn(); err != nil {
		return nil, err
	}
	rabbitMQ.conn = conn

	channel, err := rabbitMQ.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("rabbitmq create channel failed: %w", err)
	}
	rabbitMQ.channel = channel
	if err := channel.Confirm(false); err != nil {
		_ = channel.Close()
		return nil, fmt.Errorf("rabbitmq enable publisher confirms: %w", err)
	}
	rabbitMQ.confirms = channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	return rabbitMQ, nil
}

func (r *RabbitMQ) declareQueue() (amqp.Queue, error) {
	if r == nil || r.channel == nil {
		return amqp.Queue{}, errors.New("rabbitmq channel unavailable")
	}

	return r.channel.QueueDeclare(
		r.Key,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMQ) Publish(message []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, err := r.declareQueue(); err != nil {
		return err
	}

	if err := r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         message,
			DeliveryMode: amqp.Persistent,
		},
	); err != nil {
		return err
	}

	select {
	case confirmation, ok := <-r.confirms:
		if !ok {
			return errors.New("rabbitmq publisher confirmation channel closed")
		}
		if !confirmation.Ack {
			return fmt.Errorf("rabbitmq rejected published message %d", confirmation.DeliveryTag)
		}
		return nil
	case <-time.After(5 * time.Second):
		return errors.New("rabbitmq publisher confirmation timed out")
	}
}

func (r *RabbitMQ) Consume(handle func(msg *amqp.Delivery) error) {
	r.mu.Lock()
	if _, err := r.declareQueue(); err != nil {
		r.mu.Unlock()
		observe.Warn(context.Background(), "rabbitmq declare queue failed", "cause", err.Error(), "queue", r.Key)
		return
	}

	if err := r.channel.Qos(8, 0, false); err != nil {
		r.mu.Unlock()
		observe.Warn(context.Background(), "rabbitmq qos setup failed", "cause", err.Error(), "queue", r.Key)
		return
	}

	msgs, err := r.channel.Consume(r.Key, "", false, false, false, false, nil)
	r.mu.Unlock()
	if err != nil {
		observe.Warn(context.Background(), "rabbitmq consume setup failed", "cause", err.Error(), "queue", r.Key)
		return
	}

	for msg := range msgs {
		if err := handle(&msg); err != nil {
			if errors.Is(err, ErrDropMessage) {
				observe.Warn(context.Background(), "rabbitmq dropping poison message", "cause", err.Error(), "queue", r.Key)
				if rejectErr := msg.Reject(false); rejectErr != nil {
					observe.Warn(context.Background(), "rabbitmq reject poison message failed", "cause", rejectErr.Error(), "queue", r.Key)
				}
				continue
			}

			observe.Warn(context.Background(), "rabbitmq consume failed, requeueing message", "cause", err.Error(), "queue", r.Key)
			if nackErr := msg.Nack(false, true); nackErr != nil {
				observe.Warn(context.Background(), "rabbitmq nack failed", "cause", nackErr.Error(), "queue", r.Key)
			}
			continue
		}

		if err := msg.Ack(false); err != nil {
			observe.Warn(context.Background(), "rabbitmq ack failed", "cause", err.Error(), "queue", r.Key)
		}
	}
}
