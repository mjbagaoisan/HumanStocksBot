package repos

import (
	"context"
	"testing"
)

func TestWalletRepo_Create(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	memberRepo := NewGuildMemberRepo(pool)
	walletRepo := NewWalletRepo(pool)
	ctx := context.Background()

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)

	guildID := "test-guild-1"
	userID := "test-user-1"
	startingCash := int64(100000)

	memberRepo.Create(ctx, tx, guildID, userID)
	
	if err := walletRepo.Create(ctx, tx, guildID, userID, startingCash); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wallet, err := walletRepo.Get(ctx, tx, guildID, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if wallet.Cash != startingCash {
		t.Errorf("got cash %d, want %d", wallet.Cash, startingCash)
	}
}

func TestWalletRepo_Update(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	memberRepo := NewGuildMemberRepo(pool)
	walletRepo := NewWalletRepo(pool)
	ctx := context.Background()

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)

	guildID := "test-guild-2"
	userID := "test-user-2"

	memberRepo.Create(ctx, tx, guildID, userID)
	walletRepo.Create(ctx, tx, guildID, userID, 100000)

	wallet, _ := walletRepo.GetForUpdate(ctx, tx, guildID, userID)
	wallet.Cash = 50000

	if err := walletRepo.Update(ctx, tx, wallet); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, _ := walletRepo.Get(ctx, tx, guildID, userID)
	if updated.Cash != 50000 {
		t.Errorf("got cash %d, want 50000", updated.Cash)
	}
}
