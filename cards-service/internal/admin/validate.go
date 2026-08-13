package admin

import (
	"context"
	"errors"
	"fmt"
	"regexp"
)

const maxIDLength = 64

var idPattern = regexp.MustCompile(`^[A-Za-z0-9_.:-]+$`)

var (
	visibilities = []Visibility{VisibilityPrivate, VisibilityPublic, VisibilityBoth}
	predicateOps = []PredicateOp{OpGt, OpGte, OpEq, OpExists}
	matchModes   = []MatchMode{MatchAll, MatchAny}
	valueTypes   = []MetricValueType{
		MetricValueNumber, MetricValueMoney, MetricValuePercentile,
		MetricValueRatio, MetricValueString, MetricValueDate,
	}
	sourceFields = []MetricSourceField{SourceFieldValue, SourceFieldPercentile, SourceFieldShare}
	sceneTypes   = []SceneType{
		SceneIntro, SceneStat, SceneInsight, SceneAchievement,
		SceneUpsell, SceneBlocks, SceneOutro, SceneCustom,
	}
)

func allowed[T ~string](value T, options []T) bool {
	for _, option := range options {
		if value == option {
			return true
		}
	}
	return false
}

func enumMessage[T ~string](options []T) string {
	names := make([]string, 0, len(options))
	for _, option := range options {
		names = append(names, string(option))
	}
	return fmt.Sprintf("must be one of %v", names)
}

func validateID(value string) error {
	if value == "" {
		return errors.New("must not be empty")
	}
	if len(value) > maxIDLength {
		return fmt.Errorf("must be at most %d characters", maxIDLength)
	}
	if !idPattern.MatchString(value) {
		return errors.New("must contain only letters, digits and _ . : -")
	}
	return nil
}

func resolveID(urlID, bodyID string, field string) (string, error) {
	if err := validateID(urlID); err != nil {
		return "", &ValidationError{
			Message: "invalid " + field,
			Fields:  fieldErrors{field: err.Error()},
		}
	}
	if bodyID != "" && bodyID != urlID {
		return "", &ValidationError{
			Message: "id in body does not match id in path",
			Fields:  fieldErrors{field: "must match " + field + " from path"},
		}
	}
	return urlID, nil
}

type metricLookup struct {
	store Store
	cache map[string]*MetricDefinition
}

func newMetricLookup(store Store) *metricLookup {
	return &metricLookup{store: store, cache: map[string]*MetricDefinition{}}
}

func (l *metricLookup) find(ctx context.Context, key string) (*MetricDefinition, error) {
	if cached, ok := l.cache[key]; ok {
		return cached, nil
	}

	metric, err := l.store.GetMetric(ctx, key)
	if errors.Is(err, ErrNotFound) {
		l.cache[key] = nil
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	l.cache[key] = &metric
	return &metric, nil
}

func (l *metricLookup) validatePredicate(ctx context.Context, p Predicate, field string, errs fieldErrors) error {
	if p.Metric == "" {
		errs.add(field+".metric", "must not be empty")
	} else {
		metric, err := l.find(ctx, p.Metric)
		if err != nil {
			return err
		}
		if metric == nil {
			errs.add(field+".metric", "unknown metric key")
		}
	}

	if !allowed(p.Op, predicateOps) {
		errs.add(field+".op", enumMessage(predicateOps))
	}

	return nil
}

func validateMetricWrite(ctx context.Context, store Store, key string, write *MetricWrite) error {
	errs := fieldErrors{}

	if !allowed(write.ValueType, valueTypes) {
		errs.add("valueType", enumMessage(valueTypes))
	}

	currency := derefString(write.Currency)
	switch {
	case write.ValueType == MetricValueMoney && currency == "":
		errs.add("currency", "is required for money metrics")
	case write.ValueType != MetricValueMoney && currency != "":
		errs.add("currency", "is allowed only for money metrics")
	case currency != "" && currency != "RUB":
		errs.add("currency", "must be RUB")
	}

	if write.SourceField == "" {
		write.SourceField = SourceFieldValue
	}
	if !allowed(write.SourceField, sourceFields) {
		errs.add("sourceField", enumMessage(sourceFields))
	}

	if write.SourceKey == "" {
		write.SourceKey = key
	}
	if err := validateID(write.SourceKey); err != nil {
		errs.add("sourceKey", err.Error())
	}

	if write.ComparisonMinPercentile != nil {
		if *write.ComparisonMinPercentile < 0 || *write.ComparisonMinPercentile > 100 {
			errs.add("comparisonMinPercentile", "must be between 0 and 100")
		}
	}

	lookup := newMetricLookup(store)
	if percentileKey := derefString(write.PercentileKey); percentileKey != "" {
		referenced, err := lookup.find(ctx, percentileKey)
		if err != nil {
			return err
		}
		switch {
		case referenced == nil:
			errs.add("percentileKey", "unknown metric key")
		case referenced.ValueType != MetricValuePercentile:
			errs.add("percentileKey", "must reference a metric with valueType percentile")
		}
	}

	return errs.err("invalid metric definition")
}

func validateBadgeWrite(ctx context.Context, store Store, write *BadgeWrite) error {
	errs := fieldErrors{}

	if write.Title == "" {
		errs.add("title", "must not be empty")
	}
	if write.Description == "" {
		errs.add("description", "must not be empty")
	}
	if !allowed(write.Visibility, visibilities) {
		errs.add("visibility", enumMessage(visibilities))
	}

	lookup := newMetricLookup(store)
	if err := lookup.validatePredicate(ctx, write.When, "when", errs); err != nil {
		return err
	}

	return errs.err("invalid badge rule")
}

func validateStoryWrite(ctx context.Context, store Store, id string, write *StoryWrite) error {
	errs := fieldErrors{}

	if !allowed(write.SceneType, sceneTypes) {
		errs.add("sceneType", enumMessage(sceneTypes))
	}
	if !allowed(write.Visibility, visibilities) {
		errs.add("visibility", enumMessage(visibilities))
	}

	if len(write.Payload) == 0 {
		errs.add("payload", "must not be empty")
	} else {
		payloadID, _ := write.Payload["id"].(string)
		payloadType, _ := write.Payload["type"].(string)

		if payloadID == "" {
			errs.add("payload.id", "must not be empty")
		} else if payloadID != id {
			errs.add("payload.id", "must match id from path")
		}

		if payloadType == "" {
			errs.add("payload.type", "must not be empty")
		} else if payloadType != string(write.SceneType) {
			errs.add("payload.type", "must match sceneType")
		}
	}

	if write.When != nil {
		lookup := newMetricLookup(store)
		if err := lookup.validatePredicate(ctx, *write.When, "when", errs); err != nil {
			return err
		}
	}

	return errs.err("invalid story rule")
}

func validateRecommendationWrite(ctx context.Context, store Store, write *RecommendationWrite) error {
	errs := fieldErrors{}

	if write.Title == "" {
		errs.add("title", "must not be empty")
	}
	if write.Text == "" {
		errs.add("text", "must not be empty")
	}
	if write.LinkLabel == "" {
		errs.add("linkLabel", "must not be empty")
	}
	if write.Path == "" {
		errs.add("path", "must not be empty")
	}

	if len(write.When.Predicates) > 0 && !allowed(write.When.Match, matchModes) {
		errs.add("when.match", enumMessage(matchModes))
	}

	lookup := newMetricLookup(store)
	for index, predicate := range write.When.Predicates {
		field := fmt.Sprintf("when.predicates[%d]", index)
		if err := lookup.validatePredicate(ctx, predicate, field, errs); err != nil {
			return err
		}
	}

	return errs.err("invalid recommendation rule")
}
