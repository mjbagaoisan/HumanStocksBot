package repos

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mjbagaoisan/humanstocksbot/internal/domain"
)

type GuildConfigRepo interface {
	GetOrCreate(ctx context.Context, tx pgx.Tx, guildID string) (*domain.GuildConfig, error)
	Update(ctx context.Context, tx pgx.Tx, config *domain.GuildConfig) error
}

type GuildMemberRepo interface {
	Get(ctx context.Context, tx pgx.Tx, guildID, userID string) (*domain.GuildMember, error)
	Exists(ctx context.Context, tx pgx.Tx, guildID, userID string) (bool, error)
	Create(ctx context.Context, tx pgx.Tx, guildID, userID string) error
	SetOptedIn(ctx context.Context, tx pgx.Tx, guildID, userID string, optedIn bool) error
}

type WalletRepo interface {
	Create(ctx context.Context, tx pgx.Tx, guildID, userID string, startingCash int64) error
	Get(ctx context.Context, tx pgx.Tx, guildID, userID string) (*domain.Wallet, error)
	GetForUpdate(ctx context.Context, tx pgx.Tx, guildID, userID string) (*domain.Wallet, error)
	Update(ctx context.Context, tx pgx.Tx, wallet *domain.Wallet) error
}

type MarketRepo interface {
	Create(ctx context.Context, tx pgx.Tx, market *domain.Market) error
	Get(ctx context.Context, tx pgx.Tx, guildID, subjectUserID string) (*domain.Market, error)
	GetForUpdate(ctx context.Context, tx pgx.Tx, guildID, subjectUserID string) (*domain.Market, error)
	GetByStatus(ctx context.Context, tx pgx.Tx, guildID, status string) ([]*domain.Market, error)
	Update(ctx context.Context, tx pgx.Tx, market *domain.Market) error
}
