package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"cards-service/internal/cards"
	"cards-service/internal/clients"
	"cards-service/internal/models"
)

type Handler struct {
	userClient      *clients.UserClient
	analyticsClient *clients.AnalyticsClient
	shareSigningKey []byte
	shareBaseURL    string
}

func NewHandler(
	userClient *clients.UserClient,
	analyticsClient *clients.AnalyticsClient,
	shareSigningKey string,
	shareBaseURL string,
) *Handler {
	return &Handler{
		userClient:      userClient,
		analyticsClient: analyticsClient,
		shareSigningKey: []byte(shareSigningKey),
		shareBaseURL:    shareBaseURL,
	}
}

func RegisterRoutes(handler *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /api/recap/{year}/{id}", handler.GetRecap)
	mux.HandleFunc("GET /api/share/{token}", handler.GetSharedRecap)

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

	return cards.BuildRecap(profile, year, metrics, cards.BuildOptions{
		Mode:         mode,
		SigningKey:   h.shareSigningKey,
		ShareBaseURL: h.shareBaseURL,
	}), nil
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
