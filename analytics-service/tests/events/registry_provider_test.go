package events_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"analytics-service/internal/events"
)

type countingLoader struct {
	calls    int
	registry *events.Registry
	err      error
}

func (l *countingLoader) Load(context.Context) (*events.Registry, error) {
	l.calls++
	if l.err != nil {
		return nil, l.err
	}
	return l.registry, nil
}

func TestRegistryProviderShouldReuseCacheWhenTTLHasNotExpired(t *testing.T) {
	loader := &countingLoader{registry: testRegistry()}
	provider := events.NewRegistryProvider(loader, time.Hour)

	if _, err := provider.Get(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := provider.Get(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if loader.calls != 1 {
		t.Fatalf("expected a single load within TTL, got %d", loader.calls)
	}
}

func TestRegistryProviderShouldReloadWhenTTLHasExpired(t *testing.T) {
	loader := &countingLoader{registry: testRegistry()}
	provider := events.NewRegistryProvider(loader, 0)

	if _, err := provider.Get(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := provider.Get(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if loader.calls != 2 {
		t.Fatalf("expected a reload after TTL, got %d loads", loader.calls)
	}
}

func TestRegistryProviderShouldServeCachedRegistryWhenReloadFails(t *testing.T) {
	loader := &countingLoader{registry: testRegistry()}
	provider := events.NewRegistryProvider(loader, 0)

	if _, err := provider.Get(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	loader.err = errors.New("database is down")

	registry, err := provider.Get(context.Background())

	if err != nil {
		t.Fatalf("expected cached registry, got error: %v", err)
	}
	if !registry.Has("item_published") {
		t.Fatal("expected cached registry contents")
	}
}

func TestRegistryProviderShouldFailWhenFirstLoadFails(t *testing.T) {
	loader := &countingLoader{err: errors.New("database is down")}
	provider := events.NewRegistryProvider(loader, time.Hour)

	if _, err := provider.Get(context.Background()); err == nil {
		t.Fatal("expected error when there is nothing cached")
	}
}
