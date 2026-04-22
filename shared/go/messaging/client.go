package messaging

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type Publisher interface {
	Publish(ctx context.Context, exchange string, routingKey string, message []byte) error
	Close() error
}

type Subscriber interface {
	Subscribe(ctx context.Context, exchange string, routingKey string, handler MessageHandler) error
	Close() error
}

type MessageHandler func(context.Context, []byte) error

type RabbitMQClient struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	mu    sync.RWMutex
	url   string
	pool  *WorkerPool
}

type WorkerPool struct {
	workers int
	ch      chan struct{}
	mu      sync.Mutex
}

func NewRabbitMQClient(url string, workerCount int) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	if err := ch.ExchangeDeclare(
		"daedalus.events",
		amqp.ExchangeTopic,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	pool := &WorkerPool{
		workers: workerCount,
		ch:      make(chan struct{}, workerCount),
	}
	for i := 0; i < workerCount; i++ {
		pool.ch <- struct{}{}
	}

	return &RabbitMQClient{
		conn: conn,
		ch:   ch,
		url:  url,
		pool: pool,
	}, nil
}

func (rc *RabbitMQClient) Publish(ctx context.Context, exchange string, routingKey string, message []byte) error {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	if rc.ch == nil {
		return fmt.Errorf("channel is closed")
	}

	return rc.ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          message,
			Timestamp:     time.Now(),
			DeliveryMode:  amqp.Persistent,
		},
	)
}

func (rc *RabbitMQClient) Subscribe(ctx context.Context, exchange string, routingKey string, handler MessageHandler) error {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	if rc.ch == nil {
		return fmt.Errorf("channel is closed")
	}

	queueName := fmt.Sprintf("%s.%s", exchange, routingKey)
	q, err := rc.ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := rc.ch.QueueBind(q.Name, routingKey, exchange, false, nil); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	msgs, err := rc.ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to consume messages: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-msgs:
				<-rc.pool.ch

				go func(m amqp.Delivery) {
					defer func() {
						rc.pool.ch <- struct{}{}
					}()

					if err := handler(ctx, m.Body); err != nil {
						m.Nack(false, true)
					} else {
						m.Ack(false)
					}
				}(msg)
			}
		}
	}()

	return nil
}

func (rc *RabbitMQClient) Close() error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if rc.ch != nil {
		rc.ch.Close()
		rc.ch = nil
	}

	if rc.conn != nil {
		return rc.conn.Close()
	}

	return nil
}

func (rc *RabbitMQClient) HealthCheck(ctx context.Context) error {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	if rc.ch == nil {
		return fmt.Errorf("channel is not available")
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	ch, err := rc.conn.Channel()
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	ch.Close()

	return nil
}
