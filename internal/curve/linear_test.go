package curve

import (
	"errors"
	"math"
	"testing"

	"github.com/mjbagaoisan/humanstocksbot/internal/domain"
)

func TestPrice(t *testing.T) {
	checkPrice := func(t *testing.T, got, want int64) {
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	}

	t.Run("check price of next share", func(t *testing.T) {
		got := Price(1000, 100, 5)
		want := int64(1500)
		checkPrice(t, got, want)
	})

}

func TestBuyCost(t *testing.T) {
	checkBuy := func(t *testing.T, got, want int64) {
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	}
	t.Run("buy 1 share at supply 0", func(t *testing.T) {
		got, err := BuyCost(1000, 100, 0, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := int64(1000)
		checkBuy(t, got, want)
	})

	t.Run("Buy multiple shares at supply 0", func(t *testing.T) {
		got, err := BuyCost(1000, 100, 0, 3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := int64(3300)
		checkBuy(t, got, want)
	})
	t.Run("where supply isn't 0", func(t *testing.T) {
		got, err := BuyCost(1000, 100, 5, 3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := int64(4800)
		checkBuy(t, got, want)

	})

	t.Run("quantity <= 0 returns ErrInvalidQuantity", func(t *testing.T) {
		_, err := BuyCost(1000, 100, 5, 0)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidQuantity) {
			t.Fatalf("expected ErrInvalidQuantity, got %v", err)
		}
	})

	t.Run("overflow returns ErrOverflow", func(t *testing.T) {
		_, err := BuyCost(math.MaxInt64, 0, 0, 2)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrOverflow) {
			t.Fatalf("expected ErrOverflow, got %v", err)
		}
	})
}

func TestSellPayout(t *testing.T) {
	checkSell := func(t *testing.T, got, want int64) {
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	}
	t.Run("payout equals to what it cost to buy those shares", func(t *testing.T) {
		got, err := SellPayout(1000, 100, 8, 3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := int64(4800)
		checkSell(t, got, want)
	})

	t.Run("quantity less than or equal to 0 returns 0", func(t *testing.T) {
		got, err := SellPayout(1000, 100, 8, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := int64(0)
		checkSell(t, got, want)
	})

	t.Run("supply is less than quantity returns insufficient quantity error", func(t *testing.T) {
		_, err := SellPayout(1000, 100, 5, 8)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInsufficientSupply) {
			t.Fatalf("expected ErrInsufficientSupply, got %v", err)
		}
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

	// Buy 3 shares.
	cost, err := BuyCost(basePrice, slope, supply, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	reserve += cost
	supply += 3

	if reserve != expectedReserve(t, supply) {
		t.Fatalf("reserve invariant failed after buy: reserve=%d supply=%d", reserve, supply)
	}

	// Sell 2 shares.
	payout, err := SellPayout(basePrice, slope, supply, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	reserve -= payout
	supply -= 2

	if reserve != expectedReserve(t, supply) {
		t.Fatalf("reserve invariant failed after sell: reserve=%d supply=%d", reserve, supply)
	}

	// Buy 1 more share.
	cost, err = BuyCost(basePrice, slope, supply, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	reserve += cost
	supply += 1

	if reserve != expectedReserve(t, supply) {
		t.Fatalf("reserve invariant failed after second buy: reserve=%d supply=%d", reserve, supply)
	}
}
