package cards

import (
	"context"
	"sync"
	"time"
)

type RuleProvider struct {
	store    *RuleStore
	ttl      time.Duration
	fallback RuleSet

	mu      sync.Mutex
	cached  RuleSet
	loaded  bool
	expires time.Time
}

func NewRuleProvider(store *RuleStore, ttl time.Duration) *RuleProvider {
	return &RuleProvider{
		store:    store,
		ttl:      ttl,
		fallback: defaultRuleSet(),
	}
}

func (p *RuleProvider) Get(ctx context.Context) RuleSet {
	if p == nil || p.store == nil {
		return defaultRuleSet()
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	if p.loaded && now.Before(p.expires) {
		return p.cached
	}

	set, err := p.store.Load(ctx)
	if err != nil {
		if p.loaded {
			return p.cached
		}
		return p.fallback
	}

	p.cached = set
	p.loaded = true
	p.expires = now.Add(p.ttl)
	return set
}
