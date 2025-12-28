package curve

import (
	"errors"
	"math"
	"testing"

	"github.com/mjbagaoisan/humanstocksbot/internal/domain"
)

func assertEqual(t *testing.T, got, want int64) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertError(t *testing.T, err, want error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error %v, got nil", want)
	}
	if !errors.Is(err, want) {
		t.Fatalf("expected error %v, got %v", want, err)
	}
}

func TestPrice(t *testing.T) {
	tests := []struct {
		name      string
		basePrice int64
		slope     int64
		supply    int64
		want      int64
	}{
		{"price at supply 5", 1000, 100, 5, 1500},
		{"price at supply 0", 1000, 100, 0, 1000},
		{"price at supply 10", 1000, 100, 10, 2000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Price(tt.basePrice, tt.slope, tt.supply)
			assertEqual(t, got, tt.want)
		})
	}
}

func TestBuyCost(t *testing.T) {
	tests := []struct {
		name      string
		basePrice int64
		slope     int64
		supply    int64
		qty       int64
		want      int64
	}{
		{"buy 1 share at supply 0", 1000, 100, 0, 1, 1000},
		{"buy multiple shares at supply 0", 1000, 100, 0, 3, 3300},
		{"buy shares where supply isn't 0", 1000, 100, 5, 3, 4800},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuyCost(tt.basePrice, tt.slope, tt.supply, tt.qty)
			assertNoError(t, err)
			assertEqual(t, got, tt.want)
		})
	}

	t.Run("quantity <= 0 returns ErrInvalidQuantity", func(t *testing.T) {
		_, err := BuyCost(1000, 100, 5, 0)
		assertError(t, err, domain.ErrInvalidQuantity)
	})

	t.Run("overflow returns ErrOverflow", func(t *testing.T) {
		_, err := BuyCost(math.MaxInt64, 0, 0, 2)
		assertError(t, err, domain.ErrOverflow)
	})
}

func TestSellPayout(t *testing.T) {
	tests := []struct {
		name      string
		basePrice int64
		slope     int64
		supply    int64
		qty       int64
		want      int64
	}{
		{"payout equals cost to buy those shares", 1000, 100, 8, 3, 4800},
		{"quantity 0 returns 0", 1000, 100, 8, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SellPayout(tt.basePrice, tt.slope, tt.supply, tt.qty)
			assertNoError(t, err)
			assertEqual(t, got, tt.want)
		})
	}

	t.Run("supply less than quantity returns ErrInsufficientSupply", func(t *testing.T) {
		_, err := SellPayout(1000, 100, 5, 8)
		assertError(t, err, domain.ErrInsufficientSupply)
	})
}

func TestReserveInvariant(t *testing.T) {
	const (
		basePrice int64 = 1000
		slope     int64 = 100
	)

	expectedReserve := func(t *testing.T, supply int64) int64 {
		t.Helper()
		if supply == 0 {
			return 0
		}
		r, err := BuyCost(basePrice, slope, 0, supply)
		if err != nil {
			t.Fatalf("unexpected error computing expected reserve: %v", err)
		}
		return r
	}

	supply := int64(0)
	reserve := int64(0)

	cost, err := BuyCost(basePrice, slope, supply, 3)
	assertNoError(t, err)
	reserve += cost
	supply += 3
	assertEqual(t, reserve, expectedReserve(t, supply))

	payout, err := SellPayout(basePrice, slope, supply, 2)
	assertNoError(t, err)
	reserve -= payout
	supply -= 2
	assertEqual(t, reserve, expectedReserve(t, supply))

	cost, err = BuyCost(basePrice, slope, supply, 1)
	assertNoError(t, err)
	reserve += cost
	supply += 1
	assertEqual(t, reserve, expectedReserve(t, supply))
}
