package repos

import (
	"context"
	"testing"
)

func TestGuildConfigRepo_GetOrCreate(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewGuildConfigRepo(pool)
	ctx := context.Background()

	t.Run("creates new config with defaults", func(t *testing.T) {
		tx, err := pool.Begin(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		guildID := "test-guild-1"
		config, err := repo.GetOrCreate(ctx, tx, guildID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if config.GuildID != guildID {
			t.Errorf("got guild_id %s, want %s", config.GuildID, guildID)
		}
		if config.StartingCash != 100000 {
			t.Errorf("got starting_cash %d, want 100000", config.StartingCash)
		}
		if config.BasePrice != 1000 {
			t.Errorf("got base_price %d, want 1000", config.BasePrice)
		}
	})

	t.Run("returns existing config", func(t *testing.T) {
		tx, err := pool.Begin(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		guildID := "test-guild-2"
		config1, _ := repo.GetOrCreate(ctx, tx, guildID)
		config2, err := repo.GetOrCreate(ctx, tx, guildID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if config1.GuildID != config2.GuildID {
			t.Error("GetOrCreate returned different configs")
		}
	})
}

func TestGuildConfigRepo_Update(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewGuildConfigRepo(pool)
	ctx := context.Background()

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	guildID := "test-guild-3"
	config, err := repo.GetOrCreate(ctx, tx, guildID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	config.StartingCash = 50000
	config.TradeFeeBps = 300

	if err := repo.Update(ctx, tx, config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := repo.GetOrCreate(ctx, tx, guildID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.StartingCash != 50000 {
		t.Errorf("got starting_cash %d, want 50000", updated.StartingCash)
	}
	if updated.TradeFeeBps != 300 {
		t.Errorf("got trade_fee_bps %d, want 300", updated.TradeFeeBps)
	}
}
