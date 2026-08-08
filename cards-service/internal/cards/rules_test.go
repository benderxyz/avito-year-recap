package cards

import (
	"testing"

	"cards-service/internal/models"
)

func TestPredicateEvalShouldRespectOperators(t *testing.T) {
	snapshot := map[models.MetricKey]float64{
		models.MetricDealsClosed:       0,
		models.MetricListingsPublished: 5,
		models.MetricSellerRating:      4.8,
	}

	cases := []struct {
		name string
		p    predicate
		want bool
	}{
		{"gt true", predicate{models.MetricListingsPublished, opGt, 0}, true},
		{"gt false", predicate{models.MetricDealsClosed, opGt, 0}, false},
		{"gte boundary", predicate{models.MetricListingsPublished, opGte, 5}, true},
		{"eq zero", predicate{models.MetricDealsClosed, opEq, 0}, true},
		{"exists nonzero", predicate{models.MetricSellerRating, opExists, 0}, true},
		{"exists zero", predicate{models.MetricDealsClosed, opExists, 0}, false},
		{"unknown metric", predicate{models.MetricMoneyEarned, opGt, 0}, false},
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
	snapshot := map[models.MetricKey]float64{
		models.MetricListingsPublished: 5,
		models.MetricDealsClosed:       0,
		models.MetricViewsTotal:        0,
	}

	all := condition{match: matchAll, predicates: []predicate{
		{models.MetricListingsPublished, opGt, 0},
		{models.MetricDealsClosed, opEq, 0},
	}}
	if !all.eval(snapshot) {
		t.Fatal("matchAll should pass when every predicate holds")
	}

	anyMatch := condition{match: matchAny, predicates: []predicate{
		{models.MetricViewsTotal, opGt, 0},
		{models.MetricListingsPublished, opGt, 0},
	}}
	if !anyMatch.eval(snapshot) {
		t.Fatal("matchAny should pass when at least one predicate holds")
	}

	none := condition{match: matchAny, predicates: []predicate{
		{models.MetricViewsTotal, opGt, 0},
		{models.MetricDealsClosed, opGt, 0},
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

func TestMetricsSnapshotShouldCoverEnum(t *testing.T) {
	snapshot := metricsSnapshot(fullMetrics())

	keys := []models.MetricKey{
		models.MetricListingsPublished, models.MetricViewsTotal, models.MetricFavoritesReceived,
		models.MetricMessagesSent, models.MetricDealsClosed, models.MetricMoneyEarned,
		models.MetricMoneySaved, models.MetricDaysActive, models.MetricPeakDayViews,
		models.MetricSearchQueries, models.MetricCategoriesTried, models.MetricDeliveryOrders,
		models.MetricActiveListings, models.MetricSellerRating, models.MetricAvgReplySeconds,
		models.MetricFirstListingAt, models.MetricFirstDealAt,
	}

	for _, key := range keys {
		if _, ok := snapshot[key]; !ok {
			t.Fatalf("snapshot missing metric key %q", key)
		}
	}
}
