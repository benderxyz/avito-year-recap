package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"user-service/internal/users"
)

type Handler struct {
	users *users.Repository
}

func NewHandler(repo *users.Repository) *Handler {
	return &Handler{users: repo}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("PUT /users/{userID}", h.UpsertUser)
	mux.HandleFunc("GET /users/{userID}", h.GetUser)
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("user-service: OK\n"))
}

func (h *Handler) UpsertUser(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserID(r.PathValue("userID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req users.UpsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	timezone := strings.TrimSpace(req.Timezone)
	if timezone == "" {
		timezone = "UTC"
	}

	user := users.User{
		UserID:     userID,
		ExternalID: req.ExternalID,
		Username:   req.Username,
		Timezone:   timezone,
	}

	if err := h.users.Upsert(r.Context(), user); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to upsert user")
		return
	}

	saved, err := h.users.GetByID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read upserted user")
		return
	}

	writeJSON(w, http.StatusOK, saved)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserID(r.PathValue("userID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.users.GetByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, users.ErrNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get user")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func parseUserID(raw string) (uint64, error) {
	userID, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, errors.New("invalid user id")
	}
	return userID, nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
