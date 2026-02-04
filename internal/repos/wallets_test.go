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
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	guildID := "test-guild-1"
	userID := "test-user-1"
	startingCash := int64(100000)

	if err := memberRepo.Create(ctx, tx, guildID, userID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	guildID := "test-guild-2"
	userID := "test-user-2"

	if err := memberRepo.Create(ctx, tx, guildID, userID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := walletRepo.Create(ctx, tx, guildID, userID, 100000); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wallet, err := walletRepo.GetForUpdate(ctx, tx, guildID, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wallet.Cash = 50000

	if err := walletRepo.Update(ctx, tx, wallet); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := walletRepo.Get(ctx, tx, guildID, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Cash != 50000 {
		t.Errorf("got cash %d, want 50000", updated.Cash)
	}
}
