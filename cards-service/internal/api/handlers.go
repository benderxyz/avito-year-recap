package api

import (
	"encoding/json"
	"net/http"

	"cards-service/internal/cards"
	"cards-service/internal/clients"
	"cards-service/internal/models"
)

type Handler struct {
	userClient      *clients.UserClient
	analyticsClient *clients.AnalyticsClient
}

func NewHandler(userClient *clients.UserClient, analyticsClient *clients.AnalyticsClient) *Handler {
	return &Handler{
		userClient:      userClient,
		analyticsClient: analyticsClient,
	}
}

func (h *Handler) GetRecap(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if _, err := h.userClient.GetProfile(r.Context(), id); err != nil {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	metrics, err := h.analyticsClient.GetMetrics(r.Context(), id)
	if err != nil {
		http.Error(w, "metrics not found", http.StatusNotFound)
		return
	}

	recap := models.Recap{
		ProfileID: id,
		Cards:     cards.Generate(metrics),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(recap); err != nil {
		http.Error(w, "failed to encode recap", http.StatusInternalServerError)
	}
}
