package cards

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type RuleProvider struct {
	store *RuleStore
	ttl   time.Duration

	mu      sync.Mutex
	cached  RuleSet
	loaded  bool
	expires time.Time
}

func NewRuleProvider(store *RuleStore, ttl time.Duration) *RuleProvider {
	return &RuleProvider{
		store: store,
		ttl:   ttl,
	}
}

func (p *RuleProvider) Get(ctx context.Context) (RuleSet, error) {
	if p == nil || p.store == nil {
		return RuleSet{}, fmt.Errorf("rule provider is not configured")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	if p.loaded && now.Before(p.expires) {
		return p.cached, nil
	}

	set, err := p.store.Load(ctx)
	if err != nil {
		if p.loaded {
			return p.cached, nil
		}
		return RuleSet{}, fmt.Errorf("load rules: %w", err)
	}

	p.cached = set
	p.loaded = true
	p.expires = now.Add(p.ttl)
	return set, nil
}
