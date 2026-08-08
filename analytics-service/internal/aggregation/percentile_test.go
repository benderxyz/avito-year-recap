package aggregation

import "testing"

func TestLowerPercentileShouldReturnMissingWhenTotalsAreEmpty(t *testing.T) {
	result := lowerPercentile(nil, 10)

	if result.Present {
		t.Fatalf("expected missing percentile, got %+v", result)
	}
}

func TestLowerPercentileShouldCountUsersBelowValueWhenTotalsExist(t *testing.T) {
	result := lowerPercentile([]float64{1, 2, 3, 10}, 3)

	if !result.Present {
		t.Fatal("expected present percentile")
	}
	if result.Value != 50 {
		t.Fatalf("expected 50, got %v", result.Value)
	}
}

func TestHigherPercentileShouldCountUsersAboveValueWhenTotalsExist(t *testing.T) {
	result := higherPercentile([]float64{1, 2, 8, 9}, 2)

	if !result.Present {
		t.Fatal("expected present percentile")
	}
	if result.Value != 50 {
		t.Fatalf("expected 50, got %v", result.Value)
	}
}
