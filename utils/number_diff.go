package utils

import "math"

type NumberDiffResult struct {
	CashDifference    float64
	AbsCashDifference float64
	Sign              int8
	PercentDifference float64
}

func NumberDiff(originalValue, newValue float64) NumberDiffResult {
	change := newValue - originalValue
	var sign int8 = 0

	if change > 0 {
		sign = 1
	} else if change < 0 {
		sign = -1
	}

	return NumberDiffResult{
		CashDifference:    change,
		AbsCashDifference: math.Abs(change),
		Sign:              sign,
		PercentDifference: change / originalValue,
	}
}
