package repos

import (
	"context"
	"testing"
)

func TestGuildMemberRepo_Create(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewGuildMemberRepo(pool)
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

	if err := repo.Create(ctx, tx, guildID, userID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	exists, err := repo.Exists(ctx, tx, guildID, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Error("member should exist after create")
	}
}

func TestGuildMemberRepo_Get(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewGuildMemberRepo(pool)
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

	if err := repo.Create(ctx, tx, guildID, userID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	member, err := repo.Get(ctx, tx, guildID, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if member.GuildID != guildID {
		t.Errorf("got guild_id %s, want %s", member.GuildID, guildID)
	}
	if member.UserID != userID {
		t.Errorf("got user_id %s, want %s", member.UserID, userID)
	}
	if member.OptedIn {
		t.Error("new member should not be opted in by default")
	}
}

func TestGuildMemberRepo_SetOptedIn(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewGuildMemberRepo(pool)
	ctx := context.Background()

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	guildID := "test-guild-3"
	userID := "test-user-3"

	if err := repo.Create(ctx, tx, guildID, userID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := repo.SetOptedIn(ctx, tx, guildID, userID, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	member, err := repo.Get(ctx, tx, guildID, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !member.OptedIn {
		t.Error("member should be opted in")
	}
	if member.OptedInAt == nil {
		t.Error("opted_in_at should be set")
	}
}
