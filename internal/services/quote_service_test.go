package services

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/mjbagaoisan/humanstocksbot/internal/domain"
)

// --- mocks ---

type mockMarketRepo struct {
	market *domain.Market
	err    error
}

func (m *mockMarketRepo) Create(_ context.Context, _ pgx.Tx, _ *domain.Market) error { return nil }
func (m *mockMarketRepo) Get(_ context.Context, _ pgx.Tx, _, _ string) (*domain.Market, error) {
	return m.market, m.err
}
func (m *mockMarketRepo) GetForUpdate(_ context.Context, _ pgx.Tx, _, _ string) (*domain.Market, error) {
	return m.market, m.err
}
func (m *mockMarketRepo) GetByStatus(_ context.Context, _ pgx.Tx, _, _ string) ([]*domain.Market, error) {
	return nil, nil
}
func (m *mockMarketRepo) Update(_ context.Context, _ pgx.Tx, _ *domain.Market) error { return nil }

type mockConfigRepo struct {
	config *domain.GuildConfig
	err    error
}

func (m *mockConfigRepo) GetOrCreate(_ context.Context, _ pgx.Tx, _ string) (*domain.GuildConfig, error) {
	return m.config, m.err
}
func (m *mockConfigRepo) Update(_ context.Context, _ pgx.Tx, _ *domain.GuildConfig) error {
	return nil
}

// --- default test fixtures ---

func defaultConfig() *domain.GuildConfig {
	return &domain.GuildConfig{
		GuildID:       "g1",
		StartingCash:  100_000,
		BasePrice:     1000,
		Slope:         100,
		TradeFeeBps:   200,
		SubjectFeeBps: 100,
	}
}

func activeMarket(supply int64) *domain.Market {
	return &domain.Market{
		GuildID:           "g1",
		SubjectUserID:     "u1",
		Status:            domain.MarketStatusActive,
		SharesOutstanding: supply,
	}
}

// --- tests ---

func TestQuote_HappyPath(t *testing.T) {
	tests := []struct {
		name        string
		supply      int64
		side        string
		qty         int64
		wantGross   int64
		wantFee     int64
		wantNet     int64
		wantPriceB  int64
		wantPriceA  int64
	}{
		// BuyCost(1000,100,0,1) = 1000. Fee=20. Net=1020.
		{"buy 1 at supply 0", 0, "BUY", 1, 1000, 20, 1020, 1000, 1100},
		// BuyCost(1000,100,5,3) = 3000+100*(15+3) = 4800. Fee=96. Net=4896.
		{"buy 3 at supply 5", 5, "BUY", 3, 4800, 96, 4896, 1500, 1800},
		// SellPayout(1000,100,10,2) = BuyCost(1000,100,8,2) = 2000+100*(16+1) = 3700. Fee=74. Net=3626.
		{"sell 2 at supply 10", 10, "SELL", 2, 3700, 74, 3626, 2000, 1800},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewQuoteService(
				&mockMarketRepo{market: activeMarket(tt.supply)},
				&mockConfigRepo{config: defaultConfig()},
			)

			r, err := svc.Quote(context.Background(), nil, "g1", "u1", tt.side, tt.qty)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if r.GrossAmount != tt.wantGross {
				t.Errorf("GrossAmount: got %d, want %d", r.GrossAmount, tt.wantGross)
			}
			if r.TotalFee != tt.wantFee {
				t.Errorf("TotalFee: got %d, want %d", r.TotalFee, tt.wantFee)
			}
			if r.NetAmount != tt.wantNet {
				t.Errorf("NetAmount: got %d, want %d", r.NetAmount, tt.wantNet)
			}
			if r.PriceBefore != tt.wantPriceB {
				t.Errorf("PriceBefore: got %d, want %d", r.PriceBefore, tt.wantPriceB)
			}
			if r.PriceAfter != tt.wantPriceA {
				t.Errorf("PriceAfter: got %d, want %d", r.PriceAfter, tt.wantPriceA)
			}
			if r.SubjectFee+r.SystemFee != r.TotalFee {
				t.Error("SubjectFee + SystemFee != TotalFee")
			}
		})
	}
}

func TestQuote_Errors(t *testing.T) {
	tests := []struct {
		name    string
		market  *domain.Market
		side    string
		qty     int64
		wantErr error
	}{
		{"invalid side", activeMarket(0), "YOLO", 1, nil},          // non-sentinel error
		{"zero quantity", activeMarket(0), "BUY", 0, domain.ErrInvalidQuantity},
		{"closed market", func() *domain.Market {
			m := activeMarket(0)
			m.Status = domain.MarketStatusClosed
			return m
		}(), "BUY", 1, domain.ErrMarketNotActive},
		{"sunsetting blocks buy", func() *domain.Market {
			m := activeMarket(10)
			m.Status = domain.MarketStatusSunsetting
			return m
		}(), "BUY", 1, domain.ErrMarketSunsetting},
		{"sell exceeds supply", activeMarket(5), "SELL", 10, domain.ErrInsufficientSupply},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewQuoteService(
				&mockMarketRepo{market: tt.market},
				&mockConfigRepo{config: defaultConfig()},
			)

			_, err := svc.Quote(context.Background(), nil, "g1", "u1", tt.side, tt.qty)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestQuote_SunsettingAllowsSell(t *testing.T) {
	market := activeMarket(10)
	market.Status = domain.MarketStatusSunsetting

	svc := NewQuoteService(
		&mockMarketRepo{market: market},
		&mockConfigRepo{config: defaultConfig()},
	)

	_, err := svc.Quote(context.Background(), nil, "g1", "u1", "SELL", 1)
	if err != nil {
		t.Errorf("sell on sunsetting market should succeed, got %v", err)
	}
}
