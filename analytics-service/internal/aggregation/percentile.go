package aggregation

func lowerPercentile(totals []float64, userValue float64) Result {
	if len(totals) == 0 {
		return Result{Present: false}
	}

	below := 0
	for _, total := range totals {
		if total < userValue {
			below++
		}
	}

	return Result{
		Value:   float64(below) * 100.0 / float64(len(totals)),
		Present: true,
	}
}

func higherPercentile(totals []float64, userValue float64) Result {
	if len(totals) == 0 {
		return Result{Present: false}
	}

	above := 0
	for _, total := range totals {
		if total > userValue {
			above++
		}
	}

	return Result{
		Value:   float64(above) * 100.0 / float64(len(totals)),
		Present: true,
	}
}
