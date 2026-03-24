package mq

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"server/infra/config"

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

	return r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         message,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (r *RabbitMQ) Consume(handle func(msg *amqp.Delivery) error) {
	r.mu.Lock()
	if _, err := r.declareQueue(); err != nil {
		r.mu.Unlock()
		log.Printf("rabbitmq declare queue failed: %v", err)
		return
	}

	if err := r.channel.Qos(8, 0, false); err != nil {
		r.mu.Unlock()
		log.Printf("rabbitmq qos setup failed: %v", err)
		return
	}

	msgs, err := r.channel.Consume(r.Key, "", false, false, false, false, nil)
	r.mu.Unlock()
	if err != nil {
		log.Printf("rabbitmq consume setup failed: %v", err)
		return
	}

	for msg := range msgs {
		if err := handle(&msg); err != nil {
			if errors.Is(err, ErrDropMessage) {
				log.Printf("rabbitmq dropping poison message: %v", err)
				if rejectErr := msg.Reject(false); rejectErr != nil {
					log.Printf("rabbitmq reject poison message failed: %v", rejectErr)
				}
				continue
			}

			log.Printf("rabbitmq consume failed, requeueing message: %v", err)
			if nackErr := msg.Nack(false, true); nackErr != nil {
				log.Printf("rabbitmq nack failed: %v", nackErr)
			}
			continue
		}

		if err := msg.Ack(false); err != nil {
			log.Printf("rabbitmq ack failed: %v", err)
		}
	}
}
