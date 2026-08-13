package admin

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"cards-service/internal/cards"
	"cards-service/internal/clients"
	"cards-service/internal/models"
)

var previewNames = []string{"Мария", "Алексей", "Ольга", "Дмитрий", "Ирина", "Никита"}

const previewMissingChance = 0.15

func (h *Handler) Preview(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	year := time.Now().UTC().Year()
	if raw := query.Get("year"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			h.fail(w, r, &ValidationError{
				Message: "invalid query parameter",
				Fields:  fieldErrors{"year": "must be an integer"},
			})
			return
		}
		year = parsed
	}

	mode := models.RecapModePrivate
	if raw := query.Get("mode"); raw != "" {
		switch raw {
		case string(models.RecapModePrivate):
			mode = models.RecapModePrivate
		case string(models.RecapModePublic):
			mode = models.RecapModePublic
		default:
			h.fail(w, r, &ValidationError{
				Message: "invalid query parameter",
				Fields:  fieldErrors{"mode": "must be one of [private public]"},
			})
			return
		}
	}

	seed := time.Now().UnixNano()
	if raw := query.Get("seed"); raw != "" {
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			h.fail(w, r, &ValidationError{
				Message: "invalid query parameter",
				Fields:  fieldErrors{"seed": "must be an integer"},
			})
			return
		}
		seed = parsed
	}

	if h.rules == nil {
		h.fail(w, r, errors.New("rules provider is not configured"))
		return
	}

	ruleSet, err := h.rules.Get(r.Context())
	if err != nil {
		h.fail(w, r, err)
		return
	}

	definitions := ruleSet.MetricDefinitions()
	if len(definitions) == 0 {
		h.fail(w, r, ErrNotFound)
		return
	}

	//nolint:gosec
	random := rand.New(rand.NewSource(seed))

	payload := cards.BuildRecap(
		previewProfile(random),
		year,
		previewMetrics(definitions, year, random),
		cards.BuildOptions{
			Mode:           mode,
			SigningKey:     h.shareSigningKey,
			ShareBaseURL:   h.shareBaseURL,
			ProductBaseURL: h.productBaseURL,
			Rules:          &ruleSet,
		},
	)

	writeJSON(w, http.StatusOK, payload)
}

func previewProfile(random *rand.Rand) models.Profile {
	id := random.Uint64()%900000 + 100000

	return models.Profile{
		UserID:     id,
		ExternalID: "preview-" + strconv.FormatUint(id, 10),
		Username:   previewNames[random.Intn(len(previewNames))],
		Timezone:   "Europe/Moscow",
	}
}

func previewMetrics(definitions []models.MetricDefinition, year int, random *rand.Rand) clients.Metrics {
	metrics := clients.Metrics{}

	for _, definition := range definitions {
		if random.Float64() < previewMissingChance {
			continue
		}

		sourceKey := definition.SourceKey
		if sourceKey == "" {
			sourceKey = definition.Key
		}

		sample := metrics[sourceKey]
		value := previewValue(definition, year, random)

		switch definition.SourceField {
		case models.MetricSourcePercentile:
			sample.Percentile = &value
		case models.MetricSourceShare:
			sample.Share = &value
		default:
			sample.Value = &value
		}

		metrics[sourceKey] = sample
	}

	return metrics
}

func previewValue(definition models.MetricDefinition, year int, random *rand.Rand) float64 {
	if definition.SourceField == models.MetricSourcePercentile ||
		definition.ValueType == models.MetricTypePercentile {
		return float64(random.Intn(99) + 1)
	}

	switch definition.ValueType {
	case models.MetricTypeMoney:
		return float64((random.Intn(5000) + 10) * 100)
	case models.MetricTypeRatio:
		return float64(random.Intn(101)) / 100
	case models.MetricTypeDate:
		start := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()
		return float64(start + random.Int63n(365*24*60*60))
	default:
		if definition.SourceField == models.MetricSourceShare {
			return float64(random.Intn(101)) / 100
		}
		return float64(random.Intn(5000))
	}
}
