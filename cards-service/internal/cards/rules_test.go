package cards

import (
	"testing"

	"cards-service/internal/models"
)

func TestPredicateEvalShouldRespectOperators(t *testing.T) {
	snapshot := metricSnapshot{
		"dealsClosed":       0,
		"listingsPublished": 5,
		"sellerRating":      4.8,
	}

	cases := []struct {
		name string
		p    predicate
		want bool
	}{
		{"gt true", predicate{"listingsPublished", opGt, 0}, true},
		{"gt false", predicate{"dealsClosed", opGt, 0}, false},
		{"gte boundary", predicate{"listingsPublished", opGte, 5}, true},
		{"eq zero", predicate{"dealsClosed", opEq, 0}, true},
		{"exists nonzero", predicate{"sellerRating", opExists, 0}, true},
		{"exists zero", predicate{"dealsClosed", opExists, 0}, false},
		{"unknown metric", predicate{"moneyEarned", opGt, 0}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.p.eval(snapshot); got != tc.want {
				t.Fatalf("predicate %+v: want %v, got %v", tc.p, tc.want, got)
			}
		})
	}
}

func TestConditionEvalShouldHandleMatchModes(t *testing.T) {
	snapshot := metricSnapshot{
		"listingsPublished": 5,
		"dealsClosed":       0,
		"viewsTotal":        0,
	}

	all := condition{match: matchAll, predicates: []predicate{
		{"listingsPublished", opGt, 0},
		{"dealsClosed", opEq, 0},
	}}
	if !all.eval(snapshot) {
		t.Fatal("matchAll should pass when every predicate holds")
	}

	anyMatch := condition{match: matchAny, predicates: []predicate{
		{"viewsTotal", opGt, 0},
		{"listingsPublished", opGt, 0},
	}}
	if !anyMatch.eval(snapshot) {
		t.Fatal("matchAny should pass when at least one predicate holds")
	}

	none := condition{match: matchAny, predicates: []predicate{
		{"viewsTotal", opGt, 0},
		{"dealsClosed", opGt, 0},
	}}
	if none.eval(snapshot) {
		t.Fatal("matchAny should fail when no predicate holds")
	}

	empty := condition{}
	if !empty.eval(snapshot) {
		t.Fatal("empty condition should always pass")
	}
}

func TestVisibilityAllowsShouldGateByMode(t *testing.T) {
	cases := []struct {
		v       visibility
		mode    models.RecapMode
		allowed bool
	}{
		{visibilityBoth, models.RecapModePrivate, true},
		{visibilityBoth, models.RecapModePublic, true},
		{visibilityPrivate, models.RecapModePrivate, true},
		{visibilityPrivate, models.RecapModePublic, false},
		{visibilityPublic, models.RecapModePublic, true},
		{visibilityPublic, models.RecapModePrivate, false},
	}

	for _, tc := range cases {
		if got := tc.v.allows(tc.mode); got != tc.allowed {
			t.Fatalf("visibility %q mode %q: want %v, got %v", tc.v, tc.mode, tc.allowed, got)
		}
	}
}

func TestMetricsSnapshotShouldCoverEveryDefinitionWithData(t *testing.T) {
	defs := testMetricDefinitions()
	snapshot := metricsSnapshot(fullMetrics(), defs)

	for _, def := range defs {
		if _, ok := snapshot[models.MetricKey(def.Key)]; !ok {
			t.Fatalf("snapshot missing metric key %q", def.Key)
		}
	}
}

func TestMetricsSnapshotShouldReadPercentileFromSourceMetric(t *testing.T) {
	defs := []models.MetricDefinition{{
		Key:         "listingsPercentile",
		ValueType:   models.MetricTypePercentile,
		SourceKey:   "listingsPublished",
		SourceField: models.MetricSourcePercentile,
	}}

	snapshot := metricsSnapshot(fullMetrics(), defs)

	if snapshot["listingsPercentile"] != 88 {
		t.Fatalf("expected percentile 88 from source metric, got %v", snapshot["listingsPercentile"])
	}
}

func TestMetricsSnapshotShouldSkipMetricWhenSampleIsMissing(t *testing.T) {
	defs := []models.MetricDefinition{{
		Key:         "unknownMetric",
		ValueType:   models.MetricTypeNumber,
		SourceKey:   "unknownMetric",
		SourceField: models.MetricSourceValue,
	}}

	snapshot := metricsSnapshot(fullMetrics(), defs)

	if _, ok := snapshot["unknownMetric"]; ok {
		t.Fatal("did not expect a snapshot entry without analytics data")
	}
}
