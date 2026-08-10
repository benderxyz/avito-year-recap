package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"user-service/internal/users"
)

const maxJSONBodyBytes = 1 << 20

type Handler struct {
	users *users.Repository
}

func NewHandler(repo *users.Repository) *Handler {
	return &Handler{users: repo}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("PUT /users/{userID}", h.UpsertUser)
	mux.HandleFunc("GET /users", h.ListUsers)
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
	if decodeErr := decodeJSONBody(w, r, &req); decodeErr != nil {
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

	if upsertErr := h.users.Upsert(r.Context(), user); upsertErr != nil {
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

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	list, err := h.users.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list users")
		return
	}

	writeJSON(w, http.StatusOK, list)
}

func parseUserID(raw string) (uint64, error) {
	userID, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || userID == 0 {
		return 0, errors.New("invalid user id")
	}
	return userID, nil
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain a single json object")
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(append(body, '\n'))
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
