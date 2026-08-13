package admin

import "context"

type MetricFilter struct {
	Enabled      *bool
	IsPublic     *bool
	IncludeInLlm *bool
	ValueType    MetricValueType
	SourceField  MetricSourceField
	Search       string
}

type BadgeFilter struct {
	Enabled    *bool
	Visibility Visibility
	Metric     string
	Search     string
}

type StoryFilter struct {
	Enabled    *bool
	Visibility Visibility
	SceneType  SceneType
	Metric     string
	Search     string
}

type RecommendationFilter struct {
	Enabled     *bool
	Metric      string
	MinPriority *int
	Search      string
}

type Store interface {
	ListMetrics(ctx context.Context, filter MetricFilter) ([]MetricDefinition, error)
	GetMetric(ctx context.Context, key string) (MetricDefinition, error)
	CreateMetric(ctx context.Context, key string, write MetricWrite) (MetricDefinition, error)
	UpdateMetric(ctx context.Context, key string, write MetricWrite) (MetricDefinition, error)
	DeleteMetric(ctx context.Context, key string) error
	MetricReferences(ctx context.Context, key string) ([]string, error)

	ListBadges(ctx context.Context, filter BadgeFilter) ([]BadgeRule, error)
	GetBadge(ctx context.Context, id string) (BadgeRule, error)
	CreateBadge(ctx context.Context, id string, write BadgeWrite) (BadgeRule, error)
	UpdateBadge(ctx context.Context, id string, write BadgeWrite) (BadgeRule, error)
	DeleteBadge(ctx context.Context, id string) error
	BadgeReferences(ctx context.Context, id string) ([]string, error)

	ListStories(ctx context.Context, filter StoryFilter) ([]StoryRule, error)
	GetStory(ctx context.Context, id string) (StoryRule, error)
	CreateStory(ctx context.Context, id string, write StoryWrite) (StoryRule, error)
	UpdateStory(ctx context.Context, id string, write StoryWrite) (StoryRule, error)
	DeleteStory(ctx context.Context, id string) error

	ListRecommendations(ctx context.Context, filter RecommendationFilter) ([]RecommendationRule, error)
	GetRecommendation(ctx context.Context, id string) (RecommendationRule, error)
	CreateRecommendation(ctx context.Context, id string, write RecommendationWrite) (RecommendationRule, error)
	UpdateRecommendation(ctx context.Context, id string, write RecommendationWrite) (RecommendationRule, error)
	DeleteRecommendation(ctx context.Context, id string) error
}
