package repos

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjbagaoisan/humanstocksbot/internal/domain"
)

type walletRepo struct {
	pool *pgxpool.Pool
}

func NewWalletRepo(pool *pgxpool.Pool) WalletRepo {
	return &walletRepo{pool: pool}
}

func (r *walletRepo) Create(ctx context.Context, tx pgx.Tx, guildID, userID string, startingCash int64) error {
	const query = `
		INSERT INTO wallets (guild_id, user_id, cash)
		VALUES ($1, $2, $3)
		ON CONFLICT (guild_id, user_id) DO UPDATE
		SET cash = $3
	`

	_, err := tx.Exec(ctx, query, guildID, userID, startingCash)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	return nil
}

func (r *walletRepo) Get(ctx context.Context, tx pgx.Tx, guildID, userID string) (*domain.Wallet, error) {
	const query = `
		SELECT guild_id, user_id, cash
		FROM wallets
		WHERE guild_id = $1 AND user_id = $2
	`

	wallet := &domain.Wallet{}
	err := tx.QueryRow(ctx, query, guildID, userID).Scan(
		&wallet.GuildID,
		&wallet.UserID,
		&wallet.Cash,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return wallet, nil
}

func (r *walletRepo) GetForUpdate(ctx context.Context, tx pgx.Tx, guildID, userID string) (*domain.Wallet, error) {
	const query = `
		SELECT guild_id, user_id, cash
		FROM wallets
		WHERE guild_id = $1 AND user_id = $2
		FOR UPDATE
	`

	wallet := &domain.Wallet{}
	err := tx.QueryRow(ctx, query, guildID, userID).Scan(
		&wallet.GuildID,
		&wallet.UserID,
		&wallet.Cash,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet for update: %w", err)
	}

	return wallet, nil
}

func (r *walletRepo) Update(ctx context.Context, tx pgx.Tx, wallet *domain.Wallet) error {
	const query = `
		UPDATE wallets
		SET cash = $3
		WHERE guild_id = $1 AND user_id = $2
	`

	_, err := tx.Exec(ctx, query, wallet.GuildID, wallet.UserID, wallet.Cash)
	if err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	return nil
}
