package repos

import (
	"context"
	"testing"

	"github.com/mjbagaoisan/humanstocksbot/internal/domain"
)

func TestMarketRepo_Create(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	memberRepo := NewGuildMemberRepo(pool)
	marketRepo := NewMarketRepo(pool)
	ctx := context.Background()

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)

	guildID := "test-guild-1"
	userID := "test-user-1"

	memberRepo.Create(ctx, tx, guildID, userID)

	market := &domain.Market{
		GuildID:           guildID,
		SubjectUserID:     userID,
		Status:            domain.MarketStatusActive,
		SharesOutstanding: 0,
		ReserveBalance:    0,
		LastPrice:         1000,
	}

	if err := marketRepo.Create(ctx, tx, market); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	retrieved, err := marketRepo.Get(ctx, tx, guildID, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if retrieved.Status != domain.MarketStatusActive {
		t.Errorf("got status %s, want %s", retrieved.Status, domain.MarketStatusActive)
	}
}

func TestMarketRepo_Update(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	memberRepo := NewGuildMemberRepo(pool)
	marketRepo := NewMarketRepo(pool)
	ctx := context.Background()

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)

	guildID := "test-guild-2"
	userID := "test-user-2"

	memberRepo.Create(ctx, tx, guildID, userID)

	market := &domain.Market{
		GuildID:           guildID,
		SubjectUserID:     userID,
		Status:            domain.MarketStatusActive,
		SharesOutstanding: 0,
		ReserveBalance:    0,
		LastPrice:         1000,
	}
	marketRepo.Create(ctx, tx, market)

	retrieved, _ := marketRepo.GetForUpdate(ctx, tx, guildID, userID)
	retrieved.SharesOutstanding = 10
	retrieved.ReserveBalance = 10500
	retrieved.LastPrice = 1100

	if err := marketRepo.Update(ctx, tx, retrieved); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, _ := marketRepo.Get(ctx, tx, guildID, userID)
	if updated.SharesOutstanding != 10 {
		t.Errorf("got shares %d, want 10", updated.SharesOutstanding)
	}
	if updated.LastPrice != 1100 {
		t.Errorf("got last_price %d, want 1100", updated.LastPrice)
	}
}

func TestMarketRepo_GetByStatus(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	memberRepo := NewGuildMemberRepo(pool)
	marketRepo := NewMarketRepo(pool)
	ctx := context.Background()

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)

	guildID := "test-guild-3"

	for i := 1; i <= 3; i++ {
		userID := "test-user-" + string(rune('0'+i))
		memberRepo.Create(ctx, tx, guildID, userID)
		
		status := domain.MarketStatusActive
		if i == 3 {
			status = domain.MarketStatusSunsetting
		}

		market := &domain.Market{
			GuildID:           guildID,
			SubjectUserID:     userID,
			Status:            status,
			SharesOutstanding: 0,
			ReserveBalance:    0,
			LastPrice:         1000,
		}
		marketRepo.Create(ctx, tx, market)
	}

	activeMarkets, err := marketRepo.GetByStatus(ctx, tx, guildID, string(domain.MarketStatusActive))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(activeMarkets) != 2 {
		t.Errorf("got %d active markets, want 2", len(activeMarkets))
	}
}
