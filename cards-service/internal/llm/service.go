package llm

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"cards-service/internal/clients"
	"cards-service/internal/models"
)

const defaultTimeout = 4 * time.Second

type BadgeText struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type InsightText struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Result struct {
	Badges  map[string]BadgeText `json:"badges"`
	Insight *InsightText         `json:"insight,omitempty"`
}

type cacheStore interface {
	Get(ctx context.Context, userID string, year int, mode, version string) (Result, bool, error)
	Put(ctx context.Context, userID string, year int, mode, version string, result Result) error
}

type Service struct {
	provider Provider
	cache    cacheStore
	timeout  time.Duration
	version  string
}

func NewService(provider Provider, cache *Cache, timeout time.Duration, model string) *Service {
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return &Service{
		provider: provider,
		cache:    cache,
		timeout:  timeout,
		version:  promptVersion(model),
	}
}

func promptVersion(model string) string {
	sum := sha256.Sum256([]byte(systemPrompt + "\x00" + model))
	return hex.EncodeToString(sum[:])[:12]
}

type EnrichInput struct {
	ID               string
	Year             int
	Mode             models.RecapMode
	Metrics          clients.Metrics
	PublicMetricKeys map[string]bool
}

type EnrichReport struct {
	Source          string
	BadgesRewritten int
	BadgesFallback  int
	InsightAdded    bool
}

func (s *Service) Enrich(ctx context.Context, in EnrichInput, payload models.RecapPayload) (models.RecapPayload, EnrichReport, error) {
	report := EnrichReport{Source: "skip"}
	if s == nil || s.provider == nil {
		return payload, report, nil
	}

	if s.cache != nil {
		if cached, ok, err := s.cache.Get(ctx, in.ID, in.Year, string(in.Mode), s.version); err != nil {
			slog.Warn("llm cache get failed", "error", err)
		} else if ok {
			report.Source = "cache"
			s.apply(&payload, cached, &report)
			return payload, report, nil
		}
	}

	report.Source = "llm"
	result, err := s.generate(ctx, in, payload)
	if err != nil {
		return payload, report, err
	}

	s.apply(&payload, result, &report)
	if s.cache != nil {
		if err := s.cache.Put(ctx, in.ID, in.Year, string(in.Mode), s.version, result); err != nil {
			slog.Warn("llm cache put failed", "error", err)
		}
	}
	return payload, report, nil
}

type llmResponse struct {
	Badges  map[string]BadgeText `json:"badges"`
	Insight *InsightText         `json:"insight"`
}

func (s *Service) generate(ctx context.Context, in EnrichInput, payload models.RecapPayload) (Result, error) {
	cctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	messages := buildMessages(payload.Meta.User.DisplayName, payload.Badges, in.Metrics, in.PublicMetricKeys)
	raw, err := s.provider.Complete(cctx, CompletionRequest{
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   700,
	})
	if err != nil {
		return Result{}, fmt.Errorf("llm complete: %w", err)
	}

	slog.Debug("LLM raw response received", "bytes", len(raw))

	cleaned := extractJSONObject(stripCodeFence(raw))
	var resp llmResponse
	if err := json.Unmarshal([]byte(cleaned), &resp); err != nil {
		return Result{}, fmt.Errorf("parse llm response: %w", err)
	}

	return sanitize(resp, payload.Badges), nil
}

func sanitize(resp llmResponse, selected []models.Badge) Result {
	known := make(map[string]bool, len(selected))
	for _, b := range selected {
		known[b.ID] = true
	}

	out := Result{Badges: make(map[string]BadgeText)}
	for id, text := range resp.Badges {
		if !known[id] {
			continue
		}
		if valid, ok := validateBadge(text); ok {
			out.Badges[id] = valid
		}
	}

	if resp.Insight != nil {
		if valid, ok := validateInsight(*resp.Insight); ok {
			out.Insight = &valid
		}
	}
	return out
}

func (s *Service) apply(payload *models.RecapPayload, result Result, report *EnrichReport) {
	for i := range payload.Badges {
		if text, ok := result.Badges[payload.Badges[i].ID]; ok {
			payload.Badges[i].Title = text.Title
			payload.Badges[i].Description = text.Description
			report.BadgesRewritten++
		}
	}
	report.BadgesFallback = len(payload.Badges) - report.BadgesRewritten

	if result.Insight != nil {
		payload.Story = insertInsight(payload.Story, *result.Insight)
		report.InsightAdded = true
	}
}

func insertInsight(story []map[string]any, insight InsightText) []map[string]any {
	scene := map[string]any{
		"id":    "llm-insight",
		"type":  "insight",
		"title": insight.Title,
		"text":  insight.Text,
	}

	pos := 0
	for i, sc := range story {
		if t, _ := sc["type"].(string); t == "intro" {
			pos = i + 1
			break
		}
	}

	out := make([]map[string]any, 0, len(story)+1)
	out = append(out, story[:pos]...)
	out = append(out, scene)
	out = append(out, story[pos:]...)
	return out
}
