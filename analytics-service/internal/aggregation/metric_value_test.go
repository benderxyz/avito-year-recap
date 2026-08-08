package aggregation

import "testing"

func TestRoundMetricFloatShouldRoundToTwoDecimalPlaces(t *testing.T) {
	if got := roundMetricFloat(4.756); got != 4.76 {
		t.Fatalf("expected 4.76, got %v", got)
	}
	if got := roundMetricFloat(88.444); got != 88.44 {
		t.Fatalf("expected 88.44, got %v", got)
	}
	if got := roundMetricFloat(47); got != 47 {
		t.Fatalf("expected 47, got %v", got)
	}
}

func TestCalculateShareShouldRoundToTwoDecimalPlaces(t *testing.T) {
	share := calculateShare(47, 1000)
	if share == nil || *share != 4.7 {
		t.Fatalf("expected share 4.7, got %v", share)
	}

	share = calculateShare(1, 3)
	if share == nil || *share != 33.33 {
		t.Fatalf("expected share 33.33, got %v", share)
	}
}
