package aggregation

import (
	"fmt"
	"time"
)

func YearRangeInTimezone(now time.Time, timezone string) (time.Time, time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid timezone %q: %w", timezone, err)
	}

	localNow := now.In(loc)
	from := time.Date(localNow.Year(), time.January, 1, 0, 0, 0, 0, loc)

	return from.UTC(), now.UTC(), nil
}
