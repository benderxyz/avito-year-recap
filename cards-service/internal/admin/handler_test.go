package admin

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestAdminAPIShouldRejectRequestWithoutToken(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})

	recorder := request(t, handler, http.MethodGet, "/api/admin/metrics", "", nil)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("want 401 without token, got %d", recorder.Code)
	}
}

func TestAdminAPIShouldRejectWrongToken(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})

	recorder := request(t, handler, http.MethodGet, "/api/admin/metrics", "wrong-token", nil)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("want 401 for wrong token, got %d", recorder.Code)
	}
}

func TestAdminAPIShouldRejectEveryTokenWhenServiceTokenIsEmpty(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux, NewHandler(Options{Store: newMemoryStore(), Rules: &stubProvider{}}))

	recorder := request(t, mux, http.MethodGet, "/api/admin/metrics", "any-token", nil)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("want 401 when ADMIN_API_TOKEN is empty, got %d", recorder.Code)
	}
}

func TestCreateMetricShouldStoreDefinitionReadableByGet(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})

	created := request(t, handler, http.MethodPost, "/api/admin/metrics", testToken, newMetricBody("viewsTotal"))
	if created.Code != http.StatusCreated {
		t.Fatalf("want 201 on create, got %d: %s", created.Code, created.Body.String())
	}

	fetched := request(t, handler, http.MethodGet, "/api/admin/metrics/viewsTotal", testToken, nil)
	if fetched.Code != http.StatusOK {
		t.Fatalf("want 200 on read, got %d", fetched.Code)
	}

	want := decodeResponse[MetricDefinition](t, created)
	got := decodeResponse[MetricDefinition](t, fetched)
	want.CreatedAt, want.UpdatedAt = "", ""
	got.CreatedAt, got.UpdatedAt = "", ""

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %+v, got %+v", want, got)
	}
}

func TestCreateMetricShouldConflictOnExistingKey(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})
	request(t, handler, http.MethodPost, "/api/admin/metrics", testToken, newMetricBody("viewsTotal"))

	recorder := request(t, handler, http.MethodPost, "/api/admin/metrics", testToken, newMetricBody("viewsTotal"))

	if recorder.Code != http.StatusConflict {
		t.Fatalf("want 409 for duplicate key, got %d", recorder.Code)
	}
}

func TestUpdateMetricShouldReturnNotFoundForUnknownKey(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})

	recorder := request(t, handler, http.MethodPut, "/api/admin/metrics/unknown", testToken, newMetricBody("unknown"))

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("want 404 for unknown key, got %d", recorder.Code)
	}
}

func TestUpdateMetricShouldRejectKeyMismatchBetweenBodyAndPath(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})
	request(t, handler, http.MethodPost, "/api/admin/metrics", testToken, newMetricBody("viewsTotal"))

	recorder := request(t, handler, http.MethodPut, "/api/admin/metrics/viewsTotal", testToken, newMetricBody("otherKey"))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("want 400 when body key differs from path, got %d", recorder.Code)
	}
}

func TestCreateMetricShouldRejectCurrencyForNonMoneyMetric(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})
	body := newMetricBody("viewsTotal")
	currency := "RUB"
	body.Currency = &currency

	recorder := request(t, handler, http.MethodPost, "/api/admin/metrics", testToken, body)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("want 400 for currency on number metric, got %d", recorder.Code)
	}
	if field := decodeResponse[AdminError](t, recorder).Fields["currency"]; field == "" {
		t.Fatal("want currency field in error response")
	}
}

func TestListMetricsShouldIncludeDisabledDefinitions(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})
	disabled := newMetricBody("disabledMetric")
	disabled.Enabled = false
	request(t, handler, http.MethodPost, "/api/admin/metrics", testToken, disabled)

	recorder := request(t, handler, http.MethodGet, "/api/admin/metrics", testToken, nil)

	items := decodeResponse[listResponse[MetricDefinition]](t, recorder).Items
	if len(items) != 1 || items[0].Key != "disabledMetric" {
		t.Fatalf("want disabled metric in list, got %+v", items)
	}
}

func TestListMetricsShouldFilterByEnabled(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})
	disabled := newMetricBody("disabledMetric")
	disabled.Enabled = false
	request(t, handler, http.MethodPost, "/api/admin/metrics", testToken, disabled)
	request(t, handler, http.MethodPost, "/api/admin/metrics", testToken, newMetricBody("enabledMetric"))

	recorder := request(t, handler, http.MethodGet, "/api/admin/metrics?enabled=true", testToken, nil)

	items := decodeResponse[listResponse[MetricDefinition]](t, recorder).Items
	if len(items) != 1 || items[0].Key != "enabledMetric" {
		t.Fatalf("want only enabled metric, got %+v", items)
	}
}

func TestListMetricsShouldRejectUnknownValueTypeFilter(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})

	recorder := request(t, handler, http.MethodGet, "/api/admin/metrics?valueType=color", testToken, nil)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("want 400 for unknown valueType filter, got %d", recorder.Code)
	}
}

func TestDeleteMetricShouldConflictWhenBadgeReferencesIt(t *testing.T) {
	store := newMemoryStore()
	handler := newTestServer(t, store, &stubProvider{})
	request(t, handler, http.MethodPost, "/api/admin/metrics", testToken, newMetricBody("messagesSent"))
	request(t, handler, http.MethodPost, "/api/admin/badges", testToken, newBadgeBody("talkative", "messagesSent"))

	recorder := request(t, handler, http.MethodDelete, "/api/admin/metrics/messagesSent", testToken, nil)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("want 409 for referenced metric, got %d", recorder.Code)
	}
	if refs := decodeResponse[AdminError](t, recorder).References; len(refs) != 1 || refs[0] != "badge:talkative" {
		t.Fatalf("want badge reference in response, got %+v", refs)
	}
}

func TestDeleteMetricShouldRemoveUnreferencedDefinition(t *testing.T) {
	store := newMemoryStore()
	handler := newTestServer(t, store, &stubProvider{})
	request(t, handler, http.MethodPost, "/api/admin/metrics", testToken, newMetricBody("viewsTotal"))

	recorder := request(t, handler, http.MethodDelete, "/api/admin/metrics/viewsTotal", testToken, nil)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("want 204 on delete, got %d", recorder.Code)
	}
	if _, exists := store.metrics["viewsTotal"]; exists {
		t.Fatal("want metric removed from store")
	}
}

func TestCreateBadgeShouldRejectUnknownMetric(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})

	recorder := request(t, handler, http.MethodPost, "/api/admin/badges", testToken, newBadgeBody("talkative", "ghostMetric"))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("want 400 for unknown metric in condition, got %d", recorder.Code)
	}
	if field := decodeResponse[AdminError](t, recorder).Fields["when.metric"]; field == "" {
		t.Fatal("want when.metric field in error response")
	}
}

func TestCreateStoryShouldRejectPayloadIdMismatch(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})
	body := newStoryBody("intro")
	body.Payload["id"] = "another-scene"

	recorder := request(t, handler, http.MethodPost, "/api/admin/stories", testToken, body)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("want 400 when payload.id differs from id, got %d", recorder.Code)
	}
	if field := decodeResponse[AdminError](t, recorder).Fields["payload.id"]; field == "" {
		t.Fatal("want payload.id field in error response")
	}
}

func TestUpdateStoryShouldRejectPayloadTypeMismatch(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})
	request(t, handler, http.MethodPost, "/api/admin/stories", testToken, newStoryBody("intro"))
	body := newStoryBody("intro")
	body.Payload["type"] = "outro"

	recorder := request(t, handler, http.MethodPut, "/api/admin/stories/intro", testToken, body)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("want 400 when payload.type differs from sceneType, got %d", recorder.Code)
	}
}

func TestCreateRecommendationShouldRequireMatchWithPredicates(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})
	request(t, handler, http.MethodPost, "/api/admin/metrics", testToken, newMetricBody("viewsTotal"))
	body := newRecommendationBody("promote", "viewsTotal")
	body.When.Match = ""

	recorder := request(t, handler, http.MethodPost, "/api/admin/recommendations", testToken, body)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("want 400 without match, got %d", recorder.Code)
	}
}

func TestCreateRecommendationShouldAcceptEmptyPredicates(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})
	body := newRecommendationBody("promote", "")
	body.When = GroupWhen{}

	recorder := request(t, handler, http.MethodPost, "/api/admin/recommendations", testToken, body)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("want 201 for always matching recommendation, got %d: %s", recorder.Code, recorder.Body.String())
	}
}

func TestWriteRequestsShouldInvalidateRuleCache(t *testing.T) {
	provider := &stubProvider{}
	handler := newTestServer(t, newMemoryStore(), provider)

	request(t, handler, http.MethodPost, "/api/admin/metrics", testToken, newMetricBody("viewsTotal"))
	request(t, handler, http.MethodPut, "/api/admin/metrics/viewsTotal", testToken, newMetricBody("viewsTotal"))
	request(t, handler, http.MethodDelete, "/api/admin/metrics/viewsTotal", testToken, nil)

	if provider.invalidated != 3 {
		t.Fatalf("want cache invalidated after every write, got %d", provider.invalidated)
	}
}

func TestReadRequestsShouldNotInvalidateRuleCache(t *testing.T) {
	provider := &stubProvider{}
	handler := newTestServer(t, newMemoryStore(), provider)

	request(t, handler, http.MethodGet, "/api/admin/metrics", testToken, nil)

	if provider.invalidated != 0 {
		t.Fatalf("want no invalidation on read, got %d", provider.invalidated)
	}
}

func TestWriteShouldSucceedWithoutRuleProvider(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux, NewHandler(Options{Token: testToken, Store: newMemoryStore()}))

	recorder := request(t, mux, http.MethodPost, "/api/admin/metrics", testToken, newMetricBody("viewsTotal"))

	if recorder.Code != http.StatusCreated {
		t.Fatalf("want 201 without configured rule provider, got %d", recorder.Code)
	}
}

func TestPreviewShouldFailWithoutRuleProvider(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux, NewHandler(Options{Token: testToken, Store: newMemoryStore()}))

	recorder := request(t, mux, http.MethodGet, "/api/admin/preview", testToken, nil)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("want 500 without configured rule provider, got %d", recorder.Code)
	}
}

func TestOpenAPIShouldNotRequireIdInUpdateBody(t *testing.T) {
	document := openAPIDocument()

	paths, _ := document["paths"].(map[string]any)
	metric, _ := paths["/api/admin/metrics/{key}"].(map[string]any)
	put, _ := metric["put"].(map[string]any)
	body, _ := put["requestBody"].(map[string]any)
	content, _ := body["content"].(map[string]any)
	media, _ := content["application/json"].(map[string]any)
	schema, _ := media["schema"].(map[string]any)

	if ref := schema["$ref"]; ref != "#/components/schemas/MetricWrite" {
		t.Fatalf("want update body without required key, got %v", ref)
	}
}

func TestOpenAPIShouldBeServedWithoutToken(t *testing.T) {
	handler := newTestServer(t, newMemoryStore(), &stubProvider{})

	recorder := request(t, handler, http.MethodGet, "/api/admin/openapi.json", "", nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("want 200 for api description, got %d", recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), "/api/admin/metrics/{key}") {
		t.Fatal("want metrics path in openapi document")
	}
}
