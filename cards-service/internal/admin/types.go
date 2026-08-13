package admin

type Visibility string

const (
	VisibilityPrivate Visibility = "private"
	VisibilityPublic  Visibility = "public"
	VisibilityBoth    Visibility = "both"
)

type PredicateOp string

const (
	OpGt     PredicateOp = "gt"
	OpGte    PredicateOp = "gte"
	OpEq     PredicateOp = "eq"
	OpExists PredicateOp = "exists"
)

type MatchMode string

const (
	MatchAll MatchMode = "all"
	MatchAny MatchMode = "any"
)

type MetricValueType string

const (
	MetricValueNumber     MetricValueType = "number"
	MetricValueMoney      MetricValueType = "money"
	MetricValuePercentile MetricValueType = "percentile"
	MetricValueRatio      MetricValueType = "ratio"
	MetricValueString     MetricValueType = "string"
	MetricValueDate       MetricValueType = "date"
)

type MetricSourceField string

const (
	SourceFieldValue      MetricSourceField = "value"
	SourceFieldPercentile MetricSourceField = "percentile"
	SourceFieldShare      MetricSourceField = "share"
)

type SceneType string

const (
	SceneIntro       SceneType = "intro"
	SceneStat        SceneType = "stat"
	SceneInsight     SceneType = "insight"
	SceneAchievement SceneType = "achievement"
	SceneUpsell      SceneType = "upsell"
	SceneBlocks      SceneType = "blocks"
	SceneOutro       SceneType = "outro"
	SceneCustom      SceneType = "custom"
)

type Predicate struct {
	Metric string      `json:"metric"`
	Op     PredicateOp `json:"op"`
	Value  float64     `json:"value"`
}

type SimpleWhen = Predicate

type GroupWhen struct {
	Match      MatchMode   `json:"match"`
	Predicates []Predicate `json:"predicates"`
}

type MetricDefinition struct {
	Key                     string            `json:"key"`
	ValueType               MetricValueType   `json:"valueType"`
	Currency                *string           `json:"currency"`
	IsPublic                bool              `json:"isPublic"`
	PercentileKey           *string           `json:"percentileKey"`
	ComparisonMinPercentile *float64          `json:"comparisonMinPercentile"`
	SourceKey               string            `json:"sourceKey"`
	SourceField             MetricSourceField `json:"sourceField"`
	IncludeInLlm            bool              `json:"includeInLlm"`
	SortOrder               int               `json:"sortOrder"`
	Enabled                 bool              `json:"enabled"`
	CreatedAt               string            `json:"createdAt"`
	UpdatedAt               string            `json:"updatedAt"`
}

type MetricWrite struct {
	ValueType               MetricValueType   `json:"valueType"`
	Currency                *string           `json:"currency"`
	IsPublic                bool              `json:"isPublic"`
	PercentileKey           *string           `json:"percentileKey"`
	ComparisonMinPercentile *float64          `json:"comparisonMinPercentile"`
	SourceKey               string            `json:"sourceKey"`
	SourceField             MetricSourceField `json:"sourceField"`
	IncludeInLlm            bool              `json:"includeInLlm"`
	SortOrder               int               `json:"sortOrder"`
	Enabled                 bool              `json:"enabled"`
}

type MetricCreate struct {
	Key string `json:"key"`
	MetricWrite
}

type BadgeRule struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IconURL     *string    `json:"iconUrl"`
	Visibility  Visibility `json:"visibility"`
	When        SimpleWhen `json:"when"`
	SortOrder   int        `json:"sortOrder"`
	Enabled     bool       `json:"enabled"`
	CreatedAt   string     `json:"createdAt"`
	UpdatedAt   string     `json:"updatedAt"`
}

type BadgeWrite struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IconURL     *string    `json:"iconUrl"`
	Visibility  Visibility `json:"visibility"`
	When        SimpleWhen `json:"when"`
	SortOrder   int        `json:"sortOrder"`
	Enabled     bool       `json:"enabled"`
}

type BadgeCreate struct {
	ID string `json:"id"`
	BadgeWrite
}

type StoryRule struct {
	ID         string         `json:"id"`
	SceneType  SceneType      `json:"sceneType"`
	Visibility Visibility     `json:"visibility"`
	When       *SimpleWhen    `json:"when"`
	Payload    map[string]any `json:"payload"`
	SortOrder  int            `json:"sortOrder"`
	Enabled    bool           `json:"enabled"`
	CreatedAt  string         `json:"createdAt"`
	UpdatedAt  string         `json:"updatedAt"`
}

type StoryWrite struct {
	SceneType  SceneType      `json:"sceneType"`
	Visibility Visibility     `json:"visibility"`
	When       *SimpleWhen    `json:"when"`
	Payload    map[string]any `json:"payload"`
	SortOrder  int            `json:"sortOrder"`
	Enabled    bool           `json:"enabled"`
}

type StoryCreate struct {
	ID string `json:"id"`
	StoryWrite
}

type RecommendationRule struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Callout   string    `json:"callout"`
	LinkLabel string    `json:"linkLabel"`
	Path      string    `json:"path"`
	Priority  int       `json:"priority"`
	When      GroupWhen `json:"when"`
	Enabled   bool      `json:"enabled"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
}

type RecommendationWrite struct {
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Callout   string    `json:"callout"`
	LinkLabel string    `json:"linkLabel"`
	Path      string    `json:"path"`
	Priority  int       `json:"priority"`
	When      GroupWhen `json:"when"`
	Enabled   bool      `json:"enabled"`
}

type RecommendationCreate struct {
	ID string `json:"id"`
	RecommendationWrite
}

type listResponse[T any] struct {
	Items []T `json:"items"`
}

type AdminError struct {
	Error      string            `json:"error"`
	Fields     map[string]string `json:"fields,omitempty"`
	References []string          `json:"references,omitempty"`
}
