package admin

import (
	"math/rand"
	"net/http"
	"testing"

	"cards-service/internal/models"
)

func testRandom(seed int64) *rand.Rand {
	//nolint:gosec
	return rand.New(rand.NewSource(seed))
}

func previewDefinition(key string, valueType models.MetricType, sourceField models.MetricSourceField) models.MetricDefinition {
	return models.MetricDefinition{
		Key:         key,
		ValueType:   valueType,
		SourceKey:   key,
		SourceField: sourceField,
	}
}

func TestPreviewShouldRejectNonNumericYear(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})

	recorder := request(t, handler, http.MethodGet, "/api/admin/preview?year=soon", testToken, nil)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("want 400 for non numeric year, got %d", recorder.Code)
	}
}

func TestPreviewShouldRejectUnknownMode(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})

	recorder := request(t, handler, http.MethodGet, "/api/admin/preview?mode=secret", testToken, nil)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("want 400 for unknown mode, got %d", recorder.Code)
	}
}

func TestPreviewShouldReturnNotFoundWithoutMetricDefinitions(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})

	recorder := request(t, handler, http.MethodGet, "/api/admin/preview", testToken, nil)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("want 404 when admin has no metrics, got %d", recorder.Code)
	}
}

func TestPreviewShouldRequireToken(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})

	recorder := request(t, handler, http.MethodGet, "/api/admin/preview", "", nil)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("want 401 without token, got %d", recorder.Code)
	}
}

func TestPreviewMetricsShouldFillFieldReadByDefinition(t *testing.T) {
	definitions := []models.MetricDefinition{
		previewDefinition("viewsPercentile", models.MetricTypePercentile, models.MetricSourcePercentile),
	}

	metrics := previewMetrics(definitions, 2026, testRandom(1))

	sample, ok := metrics["viewsPercentile"]
	if !ok {
		t.Fatal("want sample for percentile definition")
	}
	if sample.Percentile == nil || sample.Value != nil {
		t.Fatalf("want percentile field filled only, got %+v", sample)
	}
}

func TestPreviewMetricsShouldKeepPercentileInRange(t *testing.T) {
	definitions := []models.MetricDefinition{
		previewDefinition("viewsPercentile", models.MetricTypePercentile, models.MetricSourcePercentile),
	}

	metrics := previewMetrics(definitions, 2026, testRandom(7))

	sample := metrics["viewsPercentile"]
	if sample.Percentile == nil {
		t.Fatal("want percentile value")
	}
	if *sample.Percentile < 1 || *sample.Percentile > 99 {
		t.Fatalf("want percentile between 1 and 99, got %v", *sample.Percentile)
	}
}

func TestPreviewMetricsShouldRepeatValuesForSameSeed(t *testing.T) {
	definitions := []models.MetricDefinition{
		previewDefinition("viewsTotal", models.MetricTypeNumber, models.MetricSourceValue),
		previewDefinition("moneyEarned", models.MetricTypeMoney, models.MetricSourceValue),
	}

	first := previewMetrics(definitions, 2026, testRandom(42))
	second := previewMetrics(definitions, 2026, testRandom(42))

	for key, sample := range first {
		other, ok := second[key]
		if !ok || (sample.Value == nil) != (other.Value == nil) {
			t.Fatalf("want same metrics for same seed, %q differs", key)
		}
		if sample.Value != nil && *sample.Value != *other.Value {
			t.Fatalf("want same value for %q, got %v and %v", key, *sample.Value, *other.Value)
		}
	}
}

func TestPreviewValueShouldStayInsideRequestedYear(t *testing.T) {
	definition := previewDefinition("firstListingAt", models.MetricTypeDate, models.MetricSourceValue)

	value := previewValue(definition, 2026, testRandom(3))

	const startOf2026, startOf2027 = 1767225600, 1798761600
	if value < startOf2026 || value >= startOf2027 {
		t.Fatalf("want unix timestamp inside 2026, got %v", value)
	}
}
