package curve

import (
	"math"
	"math/bits"

	"github.com/mjbagaoisan/humanstocksbot/internal/domain"
)

func safeAdd(a, b int64) (int64, error) {
	if b > 0 && a > math.MaxInt64-b {
		return 0, domain.ErrOverflow
	}
	if b < 0 && a < math.MinInt64-b {
		return 0, domain.ErrOverflow
	}
	return a + b, nil
}

func safeNeg(a int64) (int64, error) {
	if a == math.MinInt64 {
		return 0, domain.ErrOverflow
	}
	return -a, nil
}

func safeSub(a, b int64) (int64, error) {
	negB, err := safeNeg(b)
	if err != nil {
		return 0, err
	}
	return safeAdd(a, negB)
}

func absUint64(x int64) uint64 {
	if x >= 0 {
		return uint64(x)
	}
	return uint64(^x) + 1
}

func safeMul(a, b int64) (int64, error) {
	if a == 0 || b == 0 {
		return 0, nil
	}
	if a == math.MinInt64 && b == -1 {
		return 0, domain.ErrOverflow
	}
	if b == math.MinInt64 && a == -1 {
		return 0, domain.ErrOverflow
	}

	negative := (a < 0) != (b < 0)
	ua := absUint64(a)
	ub := absUint64(b)

	hi, lo := bits.Mul64(ua, ub)
	if hi != 0 {
		return 0, domain.ErrOverflow
	}

	if negative {
		if lo > (uint64(1) << 63) {
			return 0, domain.ErrOverflow
		}
		if lo == (uint64(1) << 63) {
			return math.MinInt64, nil
		}
		return -int64(lo), nil
	}

	if lo > uint64(math.MaxInt64) {
		return 0, domain.ErrOverflow
	}
	return int64(lo), nil
}
