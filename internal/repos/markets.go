package repos

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjbagaoisan/humanstocksbot/internal/domain"
)

type marketRepo struct {
	pool *pgxpool.Pool
}

func NewMarketRepo(pool *pgxpool.Pool) MarketRepo {
	return &marketRepo{pool: pool}
}

func (r *marketRepo) Create(ctx context.Context, tx pgx.Tx, market *domain.Market) error {
	const query = `
		INSERT INTO markets (guild_id, subject_user_id, status, shares_outstanding, reserve_balance, last_price, sunset_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := tx.Exec(ctx, query,
		market.GuildID,
		market.SubjectUserID,
		market.Status,
		market.SharesOutstanding,
		market.ReserveBalance,
		market.LastPrice,
		market.SunsetAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create market: %w", err)
	}

	return nil
}

func (r *marketRepo) Get(ctx context.Context, tx pgx.Tx, guildID, subjectUserID string) (*domain.Market, error) {
	const query = `
		SELECT guild_id, subject_user_id, status, shares_outstanding, reserve_balance, last_price, sunset_at, created_at, updated_at
		FROM markets
		WHERE guild_id = $1 AND subject_user_id = $2
	`

	market := &domain.Market{}
	err := tx.QueryRow(ctx, query, guildID, subjectUserID).Scan(
		&market.GuildID,
		&market.SubjectUserID,
		&market.Status,
		&market.SharesOutstanding,
		&market.ReserveBalance,
		&market.LastPrice,
		&market.SunsetAt,
		&market.CreatedAt,
		&market.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get market: %w", err)
	}

	return market, nil
}

func (r *marketRepo) GetForUpdate(ctx context.Context, tx pgx.Tx, guildID, subjectUserID string) (*domain.Market, error) {
	const query = `
		SELECT guild_id, subject_user_id, status, shares_outstanding, reserve_balance, last_price, sunset_at, created_at, updated_at
		FROM markets
		WHERE guild_id = $1 AND subject_user_id = $2
		FOR UPDATE
	`

	market := &domain.Market{}
	err := tx.QueryRow(ctx, query, guildID, subjectUserID).Scan(
		&market.GuildID,
		&market.SubjectUserID,
		&market.Status,
		&market.SharesOutstanding,
		&market.ReserveBalance,
		&market.LastPrice,
		&market.SunsetAt,
		&market.CreatedAt,
		&market.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get market for update: %w", err)
	}

	return market, nil
}

func (r *marketRepo) GetByStatus(ctx context.Context, tx pgx.Tx, guildID, status string) ([]*domain.Market, error) {
	const query = `
		SELECT guild_id, subject_user_id, status, shares_outstanding, reserve_balance, last_price, sunset_at, created_at, updated_at
		FROM markets
		WHERE guild_id = $1 AND status = $2
		ORDER BY created_at DESC
	`

	rows, err := tx.Query(ctx, query, guildID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get markets by status: %w", err)
	}
	defer rows.Close()

	var markets []*domain.Market
	for rows.Next() {
		market := &domain.Market{}
		err := rows.Scan(
			&market.GuildID,
			&market.SubjectUserID,
			&market.Status,
			&market.SharesOutstanding,
			&market.ReserveBalance,
			&market.LastPrice,
			&market.SunsetAt,
			&market.CreatedAt,
			&market.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan market: %w", err)
		}
		markets = append(markets, market)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating markets: %w", err)
	}

	return markets, nil
}

func (r *marketRepo) Update(ctx context.Context, tx pgx.Tx, market *domain.Market) error {
	const query = `
		UPDATE markets
		SET status = $3,
		    shares_outstanding = $4,
		    reserve_balance = $5,
		    last_price = $6,
		    sunset_at = $7,
		    updated_at = NOW()
		WHERE guild_id = $1 AND subject_user_id = $2
	`

	_, err := tx.Exec(ctx, query,
		market.GuildID,
		market.SubjectUserID,
		market.Status,
		market.SharesOutstanding,
		market.ReserveBalance,
		market.LastPrice,
		market.SunsetAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update market: %w", err)
	}

	return nil
}
