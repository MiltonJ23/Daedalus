package publishers

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/Daedalus/orchestrator-agent/internal/core/domain"
)

// InMemory is a TaskPublisher used in tests / dev when Redis is not available.
type InMemory struct {
	mu        sync.Mutex
	Published []domain.SubTask
}

func NewInMemory() *InMemory { return &InMemory{} }

func (p *InMemory) Publish(_ context.Context, t domain.SubTask) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Published = append(p.Published, t)
	return "memory-" + uuid.New().String(), nil
}
