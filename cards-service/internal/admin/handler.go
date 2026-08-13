package admin

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"cards-service/internal/cards"
)

type RuleProvider interface {
	Get(ctx context.Context) (cards.RuleSet, error)
	Invalidate()
}

type Options struct {
	Token           string
	Store           Store
	Rules           RuleProvider
	ShareSigningKey []byte
	ShareBaseURL    string
	ProductBaseURL  string
}

type Handler struct {
	token           string
	store           Store
	rules           RuleProvider
	shareSigningKey []byte
	shareBaseURL    string
	productBaseURL  string
}

func NewHandler(opts Options) *Handler {
	return &Handler{
		token:           opts.Token,
		store:           opts.Store,
		rules:           opts.Rules,
		shareSigningKey: opts.ShareSigningKey,
		shareBaseURL:    opts.ShareBaseURL,
		productBaseURL:  opts.ProductBaseURL,
	}
}

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /api/admin/openapi.json", h.OpenAPI)
	mux.HandleFunc("GET /api/admin/docs", h.Docs)

	protected := http.NewServeMux()

	protected.HandleFunc("GET /api/admin/metrics", h.ListMetrics)
	protected.HandleFunc("POST /api/admin/metrics", h.CreateMetric)
	protected.HandleFunc("GET /api/admin/metrics/{key}", h.GetMetric)
	protected.HandleFunc("PUT /api/admin/metrics/{key}", h.UpdateMetric)
	protected.HandleFunc("DELETE /api/admin/metrics/{key}", h.DeleteMetric)

	protected.HandleFunc("GET /api/admin/badges", h.ListBadges)
	protected.HandleFunc("POST /api/admin/badges", h.CreateBadge)
	protected.HandleFunc("GET /api/admin/badges/{id}", h.GetBadge)
	protected.HandleFunc("PUT /api/admin/badges/{id}", h.UpdateBadge)
	protected.HandleFunc("DELETE /api/admin/badges/{id}", h.DeleteBadge)

	protected.HandleFunc("GET /api/admin/stories", h.ListStories)
	protected.HandleFunc("POST /api/admin/stories", h.CreateStory)
	protected.HandleFunc("GET /api/admin/stories/{id}", h.GetStory)
	protected.HandleFunc("PUT /api/admin/stories/{id}", h.UpdateStory)
	protected.HandleFunc("DELETE /api/admin/stories/{id}", h.DeleteStory)

	protected.HandleFunc("GET /api/admin/recommendations", h.ListRecommendations)
	protected.HandleFunc("POST /api/admin/recommendations", h.CreateRecommendation)
	protected.HandleFunc("GET /api/admin/recommendations/{id}", h.GetRecommendation)
	protected.HandleFunc("PUT /api/admin/recommendations/{id}", h.UpdateRecommendation)
	protected.HandleFunc("DELETE /api/admin/recommendations/{id}", h.DeleteRecommendation)

	protected.HandleFunc("GET /api/admin/preview", h.Preview)

	mux.Handle("/api/admin/", RequireToken(h.token, protected))
}

func (h *Handler) ListMetrics(w http.ResponseWriter, r *http.Request) {
	filter, err := parseMetricFilter(r.URL.Query())
	if err != nil {
		h.fail(w, r, err)
		return
	}

	items, err := h.store.ListMetrics(r.Context(), filter)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, listResponse[MetricDefinition]{Items: items})
}

func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	metric, err := h.store.GetMetric(r.Context(), r.PathValue("key"))
	if err != nil {
		h.fail(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, metric)
}

func (h *Handler) CreateMetric(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody[MetricCreate](r)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	key, err := requireID(body.Key, "key")
	if err != nil {
		h.fail(w, r, err)
		return
	}

	if err := validateMetricWrite(r.Context(), h.store, key, &body.MetricWrite); err != nil {
		h.fail(w, r, err)
		return
	}

	metric, err := h.store.CreateMetric(r.Context(), key, body.MetricWrite)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	h.invalidate()
	writeJSON(w, http.StatusCreated, metric)
}

func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody[MetricCreate](r)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	key, err := resolveID(r.PathValue("key"), body.Key, "key")
	if err != nil {
		h.fail(w, r, err)
		return
	}

	if err := validateMetricWrite(r.Context(), h.store, key, &body.MetricWrite); err != nil {
		h.fail(w, r, err)
		return
	}

	metric, err := h.store.UpdateMetric(r.Context(), key, body.MetricWrite)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	h.invalidate()
	writeJSON(w, http.StatusOK, metric)
}

func (h *Handler) DeleteMetric(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	references, err := h.store.MetricReferences(r.Context(), key)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	if len(references) > 0 {
		h.fail(w, r, &ReferenceError{
			Message:    "metric is referenced by other rules",
			References: references,
		})
		return
	}

	if err := h.store.DeleteMetric(r.Context(), key); err != nil {
		h.fail(w, r, err)
		return
	}

	h.invalidate()
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListBadges(w http.ResponseWriter, r *http.Request) {
	filter, err := parseBadgeFilter(r.URL.Query())
	if err != nil {
		h.fail(w, r, err)
		return
	}

	items, err := h.store.ListBadges(r.Context(), filter)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, listResponse[BadgeRule]{Items: items})
}

func (h *Handler) GetBadge(w http.ResponseWriter, r *http.Request) {
	badge, err := h.store.GetBadge(r.Context(), r.PathValue("id"))
	if err != nil {
		h.fail(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, badge)
}

func (h *Handler) CreateBadge(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody[BadgeCreate](r)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	id, err := requireID(body.ID, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}

	if err := validateBadgeWrite(r.Context(), h.store, &body.BadgeWrite); err != nil {
		h.fail(w, r, err)
		return
	}

	badge, err := h.store.CreateBadge(r.Context(), id, body.BadgeWrite)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	h.invalidate()
	writeJSON(w, http.StatusCreated, badge)
}

func (h *Handler) UpdateBadge(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody[BadgeCreate](r)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	id, err := resolveID(r.PathValue("id"), body.ID, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}

	if err := validateBadgeWrite(r.Context(), h.store, &body.BadgeWrite); err != nil {
		h.fail(w, r, err)
		return
	}

	badge, err := h.store.UpdateBadge(r.Context(), id, body.BadgeWrite)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	h.invalidate()
	writeJSON(w, http.StatusOK, badge)
}

func (h *Handler) DeleteBadge(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	references, err := h.store.BadgeReferences(r.Context(), id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	if len(references) > 0 {
		h.fail(w, r, &ReferenceError{
			Message:    "badge is referenced by story scenes",
			References: references,
		})
		return
	}

	if err := h.store.DeleteBadge(r.Context(), id); err != nil {
		h.fail(w, r, err)
		return
	}

	h.invalidate()
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListStories(w http.ResponseWriter, r *http.Request) {
	filter, err := parseStoryFilter(r.URL.Query())
	if err != nil {
		h.fail(w, r, err)
		return
	}

	items, err := h.store.ListStories(r.Context(), filter)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, listResponse[StoryRule]{Items: items})
}

func (h *Handler) GetStory(w http.ResponseWriter, r *http.Request) {
	story, err := h.store.GetStory(r.Context(), r.PathValue("id"))
	if err != nil {
		h.fail(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, story)
}

func (h *Handler) CreateStory(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody[StoryCreate](r)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	id, err := requireID(body.ID, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}

	if err := validateStoryWrite(r.Context(), h.store, id, &body.StoryWrite); err != nil {
		h.fail(w, r, err)
		return
	}

	story, err := h.store.CreateStory(r.Context(), id, body.StoryWrite)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	h.invalidate()
	writeJSON(w, http.StatusCreated, story)
}

func (h *Handler) UpdateStory(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody[StoryCreate](r)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	id, err := resolveID(r.PathValue("id"), body.ID, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}

	if err := validateStoryWrite(r.Context(), h.store, id, &body.StoryWrite); err != nil {
		h.fail(w, r, err)
		return
	}

	story, err := h.store.UpdateStory(r.Context(), id, body.StoryWrite)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	h.invalidate()
	writeJSON(w, http.StatusOK, story)
}

func (h *Handler) DeleteStory(w http.ResponseWriter, r *http.Request) {
	if err := h.store.DeleteStory(r.Context(), r.PathValue("id")); err != nil {
		h.fail(w, r, err)
		return
	}

	h.invalidate()
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListRecommendations(w http.ResponseWriter, r *http.Request) {
	filter, err := parseRecommendationFilter(r.URL.Query())
	if err != nil {
		h.fail(w, r, err)
		return
	}

	items, err := h.store.ListRecommendations(r.Context(), filter)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, listResponse[RecommendationRule]{Items: items})
}

func (h *Handler) GetRecommendation(w http.ResponseWriter, r *http.Request) {
	recommendation, err := h.store.GetRecommendation(r.Context(), r.PathValue("id"))
	if err != nil {
		h.fail(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, recommendation)
}

func (h *Handler) CreateRecommendation(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody[RecommendationCreate](r)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	id, err := requireID(body.ID, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}

	if err := validateRecommendationWrite(r.Context(), h.store, &body.RecommendationWrite); err != nil {
		h.fail(w, r, err)
		return
	}

	recommendation, err := h.store.CreateRecommendation(r.Context(), id, body.RecommendationWrite)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	h.invalidate()
	writeJSON(w, http.StatusCreated, recommendation)
}

func (h *Handler) UpdateRecommendation(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody[RecommendationCreate](r)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	id, err := resolveID(r.PathValue("id"), body.ID, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}

	if err := validateRecommendationWrite(r.Context(), h.store, &body.RecommendationWrite); err != nil {
		h.fail(w, r, err)
		return
	}

	recommendation, err := h.store.UpdateRecommendation(r.Context(), id, body.RecommendationWrite)
	if err != nil {
		h.fail(w, r, err)
		return
	}

	h.invalidate()
	writeJSON(w, http.StatusOK, recommendation)
}

func (h *Handler) DeleteRecommendation(w http.ResponseWriter, r *http.Request) {
	if err := h.store.DeleteRecommendation(r.Context(), r.PathValue("id")); err != nil {
		h.fail(w, r, err)
		return
	}

	h.invalidate()
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) invalidate() {
	if h.rules == nil {
		return
	}

	h.rules.Invalidate()
}

func (h *Handler) fail(w http.ResponseWriter, r *http.Request, err error) {
	var validation *ValidationError
	if errors.As(err, &validation) {
		writeError(w, http.StatusBadRequest, AdminError{
			Error:  validation.Message,
			Fields: validation.Fields,
		})
		return
	}

	var reference *ReferenceError
	if errors.As(err, &reference) {
		writeError(w, http.StatusConflict, AdminError{
			Error:      reference.Message,
			References: reference.References,
		})
		return
	}

	switch {
	case errors.Is(err, ErrNotFound):
		writeError(w, http.StatusNotFound, AdminError{Error: "record not found"})
	case errors.Is(err, ErrAlreadyExists):
		writeError(w, http.StatusConflict, AdminError{Error: "record already exists"})
	default:
		slog.Error("admin request failed", "error", err, "method", r.Method, "path", r.URL.Path)
		writeError(w, http.StatusInternalServerError, AdminError{Error: "internal error"})
	}
}

func decodeBody[T any](r *http.Request) (T, error) {
	var body T

	if r.Body == nil {
		return body, invalidBody("request body is required")
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		return body, invalidBody("invalid json body: " + err.Error())
	}

	return body, nil
}

func requireID(value, field string) (string, error) {
	if err := validateID(value); err != nil {
		return "", &ValidationError{
			Message: "invalid " + field,
			Fields:  fieldErrors{field: err.Error()},
		}
	}

	return value, nil
}

func parseMetricFilter(query url.Values) (MetricFilter, error) {
	enabled, err := queryBool(query, "enabled")
	if err != nil {
		return MetricFilter{}, err
	}
	isPublic, err := queryBool(query, "isPublic")
	if err != nil {
		return MetricFilter{}, err
	}
	includeInLlm, err := queryBool(query, "includeInLlm")
	if err != nil {
		return MetricFilter{}, err
	}
	valueType, err := queryEnum(query, "valueType", valueTypes)
	if err != nil {
		return MetricFilter{}, err
	}
	sourceField, err := queryEnum(query, "sourceField", sourceFields)
	if err != nil {
		return MetricFilter{}, err
	}

	return MetricFilter{
		Enabled:      enabled,
		IsPublic:     isPublic,
		IncludeInLlm: includeInLlm,
		ValueType:    valueType,
		SourceField:  sourceField,
		Search:       query.Get("search"),
	}, nil
}

func parseBadgeFilter(query url.Values) (BadgeFilter, error) {
	enabled, err := queryBool(query, "enabled")
	if err != nil {
		return BadgeFilter{}, err
	}
	vis, err := queryEnum(query, "visibility", visibilities)
	if err != nil {
		return BadgeFilter{}, err
	}

	return BadgeFilter{
		Enabled:    enabled,
		Visibility: vis,
		Metric:     query.Get("metric"),
		Search:     query.Get("search"),
	}, nil
}

func parseStoryFilter(query url.Values) (StoryFilter, error) {
	enabled, err := queryBool(query, "enabled")
	if err != nil {
		return StoryFilter{}, err
	}
	vis, err := queryEnum(query, "visibility", visibilities)
	if err != nil {
		return StoryFilter{}, err
	}
	sceneType, err := queryEnum(query, "sceneType", sceneTypes)
	if err != nil {
		return StoryFilter{}, err
	}

	return StoryFilter{
		Enabled:    enabled,
		Visibility: vis,
		SceneType:  sceneType,
		Metric:     query.Get("metric"),
		Search:     query.Get("search"),
	}, nil
}

func parseRecommendationFilter(query url.Values) (RecommendationFilter, error) {
	enabled, err := queryBool(query, "enabled")
	if err != nil {
		return RecommendationFilter{}, err
	}
	minPriority, err := queryInt(query, "minPriority")
	if err != nil {
		return RecommendationFilter{}, err
	}

	return RecommendationFilter{
		Enabled:     enabled,
		Metric:      query.Get("metric"),
		MinPriority: minPriority,
		Search:      query.Get("search"),
	}, nil
}

func queryBool(query url.Values, key string) (*bool, error) {
	raw := query.Get(key)
	if raw == "" {
		return nil, nil
	}

	parsed, err := strconv.ParseBool(raw)
	if err != nil {
		return nil, &ValidationError{
			Message: "invalid query parameter",
			Fields:  fieldErrors{key: "must be true or false"},
		}
	}

	return &parsed, nil
}

func queryInt(query url.Values, key string) (*int, error) {
	raw := query.Get(key)
	if raw == "" {
		return nil, nil
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return nil, &ValidationError{
			Message: "invalid query parameter",
			Fields:  fieldErrors{key: "must be an integer"},
		}
	}

	return &parsed, nil
}

func queryEnum[T ~string](query url.Values, key string, options []T) (T, error) {
	var empty T

	raw := query.Get(key)
	if raw == "" {
		return empty, nil
	}

	value := T(raw)
	if !allowed(value, options) {
		return empty, &ValidationError{
			Message: "invalid query parameter",
			Fields:  fieldErrors{key: enumMessage(options)},
		}
	}

	return value, nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("admin response encode failed", "error", err)
	}
}

func writeError(w http.ResponseWriter, status int, payload AdminError) {
	writeJSON(w, status, payload)
}
