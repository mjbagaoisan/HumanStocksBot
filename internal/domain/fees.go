package domain

// FeeBreakdown holds the result of a fee calculation.
type FeeBreakdown struct {
	Total      int64
	SubjectFee int64
	SystemFee  int64
}

// CalcFees computes trade fees from a gross amount using basis-point rates.
// Enforces a minimum fee of 1 on any non-zero gross amount.
func CalcFees(grossAmount, tradeFeeBps, subjectFeeBps int64) FeeBreakdown {
	total := grossAmount * tradeFeeBps / 10_000
	if total == 0 && grossAmount > 0 {
		total = 1
	}

	var subjectFee int64
	if tradeFeeBps > 0 {
		subjectFee = total * subjectFeeBps / tradeFeeBps
	}

	return FeeBreakdown{
		Total:      total,
		SubjectFee: subjectFee,
		SystemFee:  total - subjectFee,
	}
}
