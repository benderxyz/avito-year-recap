package events

import "fmt"

type Registry struct {
	byType map[string]CategoryConfig
}

func NewRegistry(configs map[string]CategoryConfig) *Registry {
	copied := make(map[string]CategoryConfig, len(configs))
	for key, value := range configs {
		copied[key] = value
	}
	return &Registry{byType: copied}
}

func (r *Registry) Get(eventType string) (CategoryConfig, error) {
	cfg, ok := r.byType[eventType]
	if !ok {
		return CategoryConfig{}, fmt.Errorf("unknown event type: %s", eventType)
	}
	return cfg, nil
}

func (r *Registry) All() map[string]CategoryConfig {
	copied := make(map[string]CategoryConfig, len(r.byType))
	for key, value := range r.byType {
		copied[key] = value
	}
	return copied
}

func (r *Registry) Has(eventType string) bool {
	_, ok := r.byType[eventType]
	return ok
}
