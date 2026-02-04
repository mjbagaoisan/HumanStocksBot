package repos

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjbagaoisan/humanstocksbot/internal/domain"
)

type guildMemberRepo struct {
	pool *pgxpool.Pool
}

func NewGuildMemberRepo(pool *pgxpool.Pool) GuildMemberRepo {
	return &guildMemberRepo{pool: pool}
}

func (r *guildMemberRepo) Get(ctx context.Context, tx pgx.Tx, guildID, userID string) (*domain.GuildMember, error) {
	const query = `
		SELECT guild_id, user_id, opted_in, opted_in_at, opted_out_at
		FROM guild_members
		WHERE guild_id = $1 AND user_id = $2
	`

	member := &domain.GuildMember{}
	err := tx.QueryRow(ctx, query, guildID, userID).Scan(
		&member.GuildID,
		&member.UserID,
		&member.OptedIn,
		&member.OptedInAt,
		&member.OptedOutAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get guild member: %w", err)
	}

	return member, nil
}

func (r *guildMemberRepo) Exists(ctx context.Context, tx pgx.Tx, guildID, userID string) (bool, error) {
	const query = `
		SELECT EXISTS(SELECT 1 FROM guild_members WHERE guild_id = $1 AND user_id = $2)
	`

	var exists bool
	err := tx.QueryRow(ctx, query, guildID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check guild member existence: %w", err)
	}

	return exists, nil
}

func (r *guildMemberRepo) Create(ctx context.Context, tx pgx.Tx, guildID, userID string) error {
	const query = `
		INSERT INTO guild_members (guild_id, user_id, opted_in)
		VALUES ($1, $2, false)
		ON CONFLICT (guild_id, user_id) DO NOTHING
	`

	_, err := tx.Exec(ctx, query, guildID, userID)
	if err != nil {
		return fmt.Errorf("failed to create guild member: %w", err)
	}

	return nil
}

func (r *guildMemberRepo) SetOptedIn(ctx context.Context, tx pgx.Tx, guildID, userID string, optedIn bool) error {
	var query string
	if optedIn {
		query = `
			UPDATE guild_members
			SET opted_in = true,
			    opted_in_at = NOW(),
			    opted_out_at = NULL
			WHERE guild_id = $1 AND user_id = $2
		`
	} else {
		query = `
			UPDATE guild_members
			SET opted_in = false,
			    opted_out_at = NOW()
			WHERE guild_id = $1 AND user_id = $2
		`
	}

	_, err := tx.Exec(ctx, query, guildID, userID)
	if err != nil {
		return fmt.Errorf("failed to set opted in status: %w", err)
	}

	return nil
}
