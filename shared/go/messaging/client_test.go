package messaging

import (
	"context"
	"testing"
)

func TestNewRabbitMQClient(t *testing.T) {
	client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/", 5)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer client.Close()

	if client.ch == nil {
		t.Error("channel should not be nil")
	}

	if client.pool.workers != 5 {
		t.Errorf("workers = %d, want 5", client.pool.workers)
	}
}

func TestPublish(t *testing.T) {
	client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/", 5)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = client.Publish(ctx, "daedalus.events", "test.event", []byte(`{"key":"value"}`))
	if err != nil {
		t.Errorf("Publish failed: %v", err)
	}
}

func TestHealthCheck(t *testing.T) {
	client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/", 5)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = client.HealthCheck(ctx)
	if err != nil {
		t.Errorf("HealthCheck failed: %v", err)
	}
}

func TestClose(t *testing.T) {
	client, err := NewRabbitMQClient("amqp://guest:guest@localhost:5672/", 5)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}

	err = client.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	if client.ch != nil {
		t.Error("channel should be nil after close")
	}
}
