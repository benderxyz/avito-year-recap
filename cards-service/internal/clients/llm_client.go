package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/responses",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("LLM API returned status %d", resp.StatusCode)
	}

	var result responseResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Output) == 0 ||
		len(result.Output[0].Content) == 0 {
		return "", fmt.Errorf("LLM response contains no text")
	}

	return result.Output[0].Content[0].Text, nil
}
