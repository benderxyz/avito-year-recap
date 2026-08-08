package aggregation

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"analytics-service/internal/events"
)

const defaultGlobalCacheTTL = 15 * time.Minute

type globalStats struct {
	db    FloatQuerier
	cache *ttlCache
}

func newGlobalStats(db FloatQuerier, ttl time.Duration) *globalStats {
	if ttl <= 0 {
		ttl = defaultGlobalCacheTTL
	}
	return &globalStats{
		db:    db,
		cache: newTTLCache(ttl),
	}
}

func (g *globalStats) counterTotal(ctx context.Context, req GlobalRequest) (Result, error) {
	key := joinCacheKey("counter_total", req.EventType, formatCacheTime(req.From), formatCacheTime(req.To))
	if cached, ok := g.cache.Get(key); ok {
		result, ok := cached.(Result)
		if !ok {
			return Result{}, fmt.Errorf("counter total cache: unexpected type %T", cached)
		}
		return result, nil
	}

	result, err := aggregateGlobalCounter(ctx, g.db, req)
	if err != nil {
		return Result{}, err
	}
	g.cache.Set(key, result)
	return result, nil
}

func (g *globalStats) uniqueTotal(
	ctx context.Context,
	req GlobalRequest,
	cfg events.CategoryConfig,
) (Result, error) {
	key := joinCacheKey(
		"unique_total",
		req.EventType,
		string(cfg.UniqueMode),
		sanitizePayloadField(cfg.UniqueField),
		req.Timezone,
		formatCacheTime(req.From),
		formatCacheTime(req.To),
	)
	if cached, ok := g.cache.Get(key); ok {
		result, ok := cached.(Result)
		if !ok {
			return Result{}, fmt.Errorf("unique total cache: unexpected type %T", cached)
		}
		return result, nil
	}

	result, err := aggregateGlobalUnique(ctx, g.db, req, cfg)
	if err != nil {
		return Result{}, err
	}
	g.cache.Set(key, result)
	return result, nil
}

func (g *globalStats) perUserTotals(
	ctx context.Context,
	key string,
	subquery string,
	args []any,
) ([]float64, error) {
	if cached, ok := g.cache.Get(key); ok {
		totals, ok := cached.([]float64)
		if !ok {
			return nil, fmt.Errorf("per-user totals cache: unexpected type %T", cached)
		}
		return totals, nil
	}

	query := fmt.Sprintf(`
		SELECT total
		FROM (%s)
	`, subquery)
	totals, err := g.db.QueryFloat64s(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("per-user totals: %w", err)
	}
	g.cache.Set(key, totals)
	return totals, nil
}

func (g *globalStats) counterTotals(ctx context.Context, req GlobalRequest) ([]float64, error) {
	subquery, args := counterPerUserSubquery(req)
	key := joinCacheKey("counter_totals", req.EventType, formatCacheTime(req.From), formatCacheTime(req.To))
	return g.perUserTotals(ctx, key, subquery, args)
}

func (g *globalStats) uniqueTotals(
	ctx context.Context,
	req GlobalRequest,
	cfg events.CategoryConfig,
) ([]float64, error) {
	subquery, args, err := uniquePerUserTotalsSubquery(cfg, req)
	if err != nil {
		return nil, err
	}
	key := joinCacheKey(
		"unique_totals",
		req.EventType,
		string(cfg.UniqueMode),
		sanitizePayloadField(cfg.UniqueField),
		req.Timezone,
		formatCacheTime(req.From),
		formatCacheTime(req.To),
	)
	return g.perUserTotals(ctx, key, subquery, args)
}

func (g *globalStats) gaugeTotals(ctx context.Context, req GlobalRequest) ([]float64, error) {
	subquery, args := gaugePerUserSubquery(req)
	key := joinCacheKey("gauge_totals", req.EventType, formatCacheTime(req.From), formatCacheTime(req.To))
	return g.perUserTotals(ctx, key, subquery, args)
}

func (g *globalStats) milestoneTotals(ctx context.Context, req GlobalRequest) ([]float64, error) {
	subquery, args := milestonePerUserSubquery(req)
	key := joinCacheKey("milestone_totals", req.EventType, formatCacheTime(req.From), formatCacheTime(req.To))
	return g.perUserTotals(ctx, key, subquery, args)
}

func (g *globalStats) intervalTotals(ctx context.Context, req GlobalRequest) ([]float64, error) {
	subquery, args := intervalPerUserSubquery(req)
	key := joinCacheKey("interval_totals", req.EventType, formatCacheTime(req.From), formatCacheTime(req.To))
	return g.perUserTotals(ctx, key, subquery, args)
}

func joinCacheKey(parts ...string) string {
	return strings.Join(parts, "|")
}

func formatCacheTime(value time.Time) string {
	return strconv.FormatInt(value.UnixNano(), 10)
}
