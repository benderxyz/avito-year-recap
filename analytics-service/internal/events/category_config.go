package events

type UniqueMode string

const (
	UniqueModePayload UniqueMode = "payload"
	UniqueModeDay     UniqueMode = "day"
)

type CategoryConfig struct {
	Category    EventCategory
	MetricKey   string
	UniqueField string
	UniqueMode  UniqueMode
}
