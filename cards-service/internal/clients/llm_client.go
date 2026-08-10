package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type LLMClient struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewLLMClient(
	apiKey string,
	baseURL string,
	model string,
	httpClient *http.Client,
) *LLMClient {
	return &LLMClient{
		apiKey:     apiKey,
		baseURL:    baseURL,
		model:      model,
		httpClient: httpClient,
	}
}

type responseRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type responseResponse struct {
	Output []struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}

const maxLLMAttempts = 3

func (c *LLMClient) Generate(
	ctx context.Context,
	input string,
) (string, error) {
	payload := responseRequest{
		Model: c.model,
		Input: input,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	var lastErr error
	for attempt := 1; attempt <= maxLLMAttempts; attempt++ {
		if attempt > 1 {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(time.Duration(attempt-1) * 300 * time.Millisecond):
			}
		}

		text, retryable, err := c.doGenerate(ctx, body)
		if err == nil {
			return text, nil
		}
		lastErr = err
		if !retryable {
			return "", err
		}
		slog.Debug("llm request retryable failure", "attempt", attempt, "error", err)
	}

	return "", fmt.Errorf("llm request failed after %d attempts: %w", maxLLMAttempts, lastErr)
}

func (c *LLMClient) doGenerate(ctx context.Context, body []byte) (string, bool, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/responses",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", false, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", ctx.Err() == nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		retryable := resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500
		return "", retryable, fmt.Errorf("LLM API returned status %d", resp.StatusCode)
	}

	var result responseResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", false, err
	}

	if text := firstNonEmptyText(result); text != "" {
		return text, false, nil
	}
	return "", false, fmt.Errorf("LLM response contains no text")
}

func firstNonEmptyText(result responseResponse) string {
	for _, out := range result.Output {
		for _, content := range out.Content {
			if content.Text != "" {
				return content.Text
			}
		}
	}
	return ""
}
