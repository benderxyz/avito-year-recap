package llm

import "context"

const (
	RoleSystem = "system"
	RoleUser   = "user"
)

type Message struct {
	Role    string
	Content string
}

type CompletionRequest struct {
	Messages    []Message
	Temperature float64
	MaxTokens   int
}

type Provider interface {
	Name() string
	Complete(ctx context.Context, req CompletionRequest) (string, error)
}
