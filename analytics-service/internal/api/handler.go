package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"analytics-service/internal/aggregation"
	"analytics-service/internal/apperr"
	"analytics-service/internal/events"
)

type Handler struct {
	ingester  *events.Ingester
	metrics   *aggregation.Service
	timezones aggregation.TimezoneResolver
}

func NewHandler(
	ingester *events.Ingester,
	metrics *aggregation.Service,
	timezones aggregation.TimezoneResolver,
) *Handler {
	return &Handler{ingester: ingester, metrics: metrics, timezones: timezones}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /events", h.IngestEvents)
	mux.HandleFunc("GET /users/{userID}/metrics", h.GetMetrics)
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("analytics-service: OK\n"))
}

func (h *Handler) IngestEvents(w http.ResponseWriter, r *http.Request) {
	var items []events.IngestEvent

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var raw json.RawMessage
	if err := decoder.Decode(&raw); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if len(raw) > 0 && raw[0] == '[' {
		if err := json.Unmarshal(raw, &items); err != nil {
			writeError(w, http.StatusBadRequest, "invalid events batch")
			return
		}
	} else {
		var single events.IngestEvent
		if err := json.Unmarshal(raw, &single); err != nil {
			writeError(w, http.StatusBadRequest, "invalid event")
			return
		}
		items = []events.IngestEvent{single}
	}

	if err := h.ingester.Ingest(r.Context(), items); err != nil {
		writeMappedError(w, err, "failed to ingest events")
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]any{
		"accepted": len(items),
	})
}

func (h *Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserID(r.PathValue("userID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	timezone := "UTC"
	if h.timezones != nil {
		resolved, err := h.timezones.Timezone(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to resolve timezone")
			return
		}
		timezone = resolved
	}

	from, to, err := ParseRange(r, timezone)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	snapshot, err := h.metrics.Metrics(r.Context(), userID, from, to)
	if err != nil {
		writeMappedError(w, err, "failed to aggregate metrics")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user_id":  userID,
		"timezone": snapshot.Timezone,
		"from":     snapshot.From.UTC().Format(time.RFC3339Nano),
		"to":       snapshot.To.UTC().Format(time.RFC3339Nano),
		"metrics":  snapshot.Metrics,
	})
}

func parseUserID(raw string) (uint64, error) {
	userID, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, errors.New("invalid user id")
	}
	return userID, nil
}

func ParseRange(r *http.Request, timezone string) (time.Time, time.Time, error) {
	fromRaw := r.URL.Query().Get("from")
	toRaw := r.URL.Query().Get("to")

	fromEmpty := fromRaw == ""
	toEmpty := toRaw == ""

	if fromEmpty != toEmpty {
		return time.Time{}, time.Time{}, errors.New("both from and to are required")
	}

	if fromEmpty && toEmpty {
		from, to, err := aggregation.YearRangeInTimezone(time.Now().UTC(), timezone)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New("invalid user timezone")
		}
		return from, to, nil
	}

	from, err := time.Parse(time.RFC3339, fromRaw)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("invalid from: use RFC3339")
	}

	to, err := time.Parse(time.RFC3339, toRaw)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("invalid to: use RFC3339")
	}

	return from.UTC(), to.UTC(), nil
}

func writeMappedError(w http.ResponseWriter, err error, internalMessage string) {
	if apperr.IsValidation(err) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeError(w, http.StatusInternalServerError, internalMessage)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		return
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
