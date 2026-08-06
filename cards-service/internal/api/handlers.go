package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"cards-service/internal/cards"
	"cards-service/internal/clients"
	"cards-service/internal/models"
	"log"
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

func RegisterRoutes(handler *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)

	mux.HandleFunc("GET /api/recap/{year}/{id}", handler.GetRecap)

	return mux
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	log.Println("Пришел запрос на проверку активность сервиса")
	if _, err := fmt.Fprintln(w, "cards-service: OK"); err != nil {
		http.Error(w, "failed to write health response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetRecap(w http.ResponseWriter, r *http.Request) {
	log.Println("Пришел запрос на получение рекапа")
	year := r.PathValue("year")
	id := r.PathValue("id")

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		http.Error(w, "invalid year", http.StatusBadRequest)
		return
	}

	// временно убрала проверку существования профиля, так как сервис user-service пока не интегрирован в cards-service
	_, err = h.userClient.GetProfile(r.Context(), id)
	log.Println("Получили профиль")
	if err != nil {
		log.Println("Ошибка при получении профиля:%w", err)
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	metrics, err := h.analyticsClient.GetMetrics(r.Context(), id, yearInt)
	if err != nil {
		http.Error(w, "metrics not found", http.StatusNotFound)
		log.Println("Ошибка при получении метрик от сервиса аналитики, err:", err)
		return
	}
	log.Printf("Метрики для пользователя %s за %d год: %+v", id, yearInt, metrics)

	recap := models.Recap{
		ProfileID: id,
		Badges:    cards.Generate(metrics),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(recap); err != nil {
		http.Error(w, "failed to encode recap", http.StatusInternalServerError)
	}
}
