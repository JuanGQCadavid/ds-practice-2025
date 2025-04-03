package core

import (
	"sync"

	"github.com/JuanGQCadavid/ds-practice-2025/order_service/internal/core/domain"
)

type Kafki struct {
	mu     sync.Mutex
	subs   []chan domain.StatesOfDemocracy
	quit   chan struct{}
	closed bool
}

func NewAgent() *Kafki {
	return &Kafki{
		subs: make([]chan domain.StatesOfDemocracy, 0),
		quit: make(chan struct{}),
	}
}

func (b *Kafki) Publish(state domain.StatesOfDemocracy) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}
	// THis could be abottle neck
	for _, ch := range b.subs {
		ch <- state
	}
}

func (b *Kafki) Subscribe() <-chan domain.StatesOfDemocracy {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil
	}

	ch := make(chan domain.StatesOfDemocracy, 100)
	b.subs = append(b.subs, ch)
	return ch
}
