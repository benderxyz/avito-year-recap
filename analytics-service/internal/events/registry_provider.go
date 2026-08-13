package events

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type registryLoader interface {
	Load(ctx context.Context) (*Registry, error)
}

type RegistryProvider struct {
	store registryLoader
	ttl   time.Duration

	mu      sync.Mutex
	cached  *Registry
	expires time.Time
}

func NewRegistryProvider(store registryLoader, ttl time.Duration) *RegistryProvider {
	return &RegistryProvider{
		store: store,
		ttl:   ttl,
	}
}

func (p *RegistryProvider) Get(ctx context.Context) (*Registry, error) {
	if p == nil || p.store == nil {
		return nil, fmt.Errorf("registry provider is not configured")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	if p.cached != nil && now.Before(p.expires) {
		return p.cached, nil
	}

	registry, err := p.store.Load(ctx)
	if err != nil {
		if p.cached != nil {
			return p.cached, nil
		}
		return nil, fmt.Errorf("load event registry: %w", err)
	}

	p.cached = registry
	p.expires = now.Add(p.ttl)
	return registry, nil
}
