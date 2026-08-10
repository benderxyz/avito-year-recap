package llm

import (
	"context"
	"strings"

	"cards-service/internal/clients"
)

type OpenAIProvider struct {
	client *clients.LLMClient
}

func NewOpenAIProvider(client *clients.LLMClient) *OpenAIProvider {
	return &OpenAIProvider{client: client}
}

func (p *OpenAIProvider) Name() string { return "openai" }

func (p *OpenAIProvider) Complete(ctx context.Context, req CompletionRequest) (string, error) {
	var b strings.Builder
	for i, m := range req.Messages {
		if i > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(m.Content)
	}
	return p.client.Generate(ctx, b.String())
}
