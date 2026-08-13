package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"cards-service/internal/admin"
	"cards-service/internal/cards"
	"cards-service/internal/clients"
	"cards-service/internal/llm"
	"cards-service/internal/models"
)

type Handler struct {
	userClient      *clients.UserClient
	analyticsClient *clients.AnalyticsClient
	shareSigningKey []byte
	shareBaseURL    string
	productBaseURL  string
	rules           *cards.RuleProvider
	llm             *llm.Service
}

func NewHandler(
	userClient *clients.UserClient,
	analyticsClient *clients.AnalyticsClient,
	shareSigningKey string,
	shareBaseURL string,
	productBaseURL string,
	rules *cards.RuleProvider,
) *Handler {
	return &Handler{
		userClient:      userClient,
		analyticsClient: analyticsClient,
		shareSigningKey: []byte(shareSigningKey),
		shareBaseURL:    shareBaseURL,
		productBaseURL:  productBaseURL,
		rules:           rules,
	}
}

func (h *Handler) SetLLMService(s *llm.Service) {
	h.llm = s
}

func RegisterRoutes(handler *Handler, adminHandler *admin.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /api/recap/{year}/{id}", handler.GetRecap)
	mux.HandleFunc("GET /api/share/{token}", handler.GetSharedRecap)

	if adminHandler != nil {
		admin.RegisterRoutes(mux, adminHandler)
	}

	return mux
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	if _, err := fmt.Fprintln(w, "cards-service: OK"); err != nil {
		http.Error(w, "failed to write health response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetRecap(w http.ResponseWriter, r *http.Request) {
	year := r.PathValue("year")
	id := r.PathValue("id")

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		http.Error(w, "invalid year", http.StatusBadRequest)
		return
	}

	recap, err := h.buildRecap(r, id, yearInt, models.RecapModePrivate)
	if err != nil {
		writeRecapError(w, err)
		return
	}

	writeJSON(w, recap)
}

func (h *Handler) GetSharedRecap(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	id, year, err := cards.DecodeShareToken(h.shareSigningKey, token)
	if err != nil {
		http.Error(w, "invalid share token", http.StatusBadRequest)
		return
	}

	recap, err := h.buildRecap(r, id, year, models.RecapModePublic)
	if err != nil {
		writeRecapError(w, err)
		return
	}

	writeJSON(w, recap)
}

func (h *Handler) buildRecap(
	r *http.Request,
	id string,
	year int,
	mode models.RecapMode,
) (models.RecapPayload, error) {
	profile, err := h.userClient.GetProfile(r.Context(), id)
	if err != nil {
		return models.RecapPayload{}, fmt.Errorf("profile not found: %w", err)
	}

	metrics, err := h.analyticsClient.GetMetrics(r.Context(), id, year, profile.Timezone)
	if err != nil {
		return models.RecapPayload{}, fmt.Errorf("metrics not found: %w", err)
	}

	if h.rules == nil {
		return models.RecapPayload{}, fmt.Errorf("rules provider is not configured")
	}

	ruleSet, err := h.rules.Get(r.Context())
	if err != nil {
		return models.RecapPayload{}, fmt.Errorf("load rules: %w", err)
	}

	opts := cards.BuildOptions{
		Mode:           mode,
		SigningKey:     h.shareSigningKey,
		ShareBaseURL:   h.shareBaseURL,
		ProductBaseURL: h.productBaseURL,
		Rules:          &ruleSet,
	}

	payload := cards.BuildRecap(profile, year, metrics, opts)

	if h.llm != nil {
		enriched, report, err := h.llm.Enrich(r.Context(), llm.EnrichInput{
			ID:          id,
			Year:        year,
			Mode:        mode,
			Metrics:     metrics,
			Definitions: ruleSet.MetricDefinitions(),
		}, payload)
		if err != nil {
			slog.Warn("llm enrich failed, serving base recap", "error", err, "user_id", id, "year", year)
		}
		slog.Debug("llm enrich",
			"source", report.Source,
			"badges_rewritten", report.BadgesRewritten,
			"badges_fallback", report.BadgesFallback,
			"insight_added", report.InsightAdded,
		)
		payload = enriched
	}

	return payload, nil
}

func writeRecapError(w http.ResponseWriter, _ error) {
	http.Error(w, "recap not found", http.StatusNotFound)
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to encode recap", http.StatusInternalServerError)
	}
}
