package repos

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjbagaoisan/humanstocksbot/internal/domain"
)

type guildConfigRepo struct {
	pool *pgxpool.Pool
}

func NewGuildConfigRepo(pool *pgxpool.Pool) GuildConfigRepo {
	return &guildConfigRepo{pool: pool}
}

func (r *guildConfigRepo) GetOrCreate(ctx context.Context, tx pgx.Tx, guildID string) (*domain.GuildConfig, error) {
	const query = `
		INSERT INTO guild_config (guild_id)
		VALUES ($1)
		ON CONFLICT (guild_id) DO NOTHING
		RETURNING guild_id, starting_cash, base_price, slope, trade_fee_bps, subject_fee_bps, sunset_days, trading_enabled
	`

	config := &domain.GuildConfig{}
	err := tx.QueryRow(ctx, query, guildID).Scan(
		&config.GuildID,
		&config.StartingCash,
		&config.BasePrice,
		&config.Slope,
		&config.TradeFeeBps,
		&config.SubjectFeeBps,
		&config.SunsetDays,
		&config.TradingPaused,
	)

	if err == pgx.ErrNoRows {
		const selectQuery = `
			SELECT guild_id, starting_cash, base_price, slope, trade_fee_bps, subject_fee_bps, sunset_days, trading_enabled
			FROM guild_config
			WHERE guild_id = $1
		`
		err = tx.QueryRow(ctx, selectQuery, guildID).Scan(
			&config.GuildID,
			&config.StartingCash,
			&config.BasePrice,
			&config.Slope,
			&config.TradeFeeBps,
			&config.SubjectFeeBps,
			&config.SunsetDays,
			&config.TradingPaused,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get or create guild config: %w", err)
	}

	config.TradingPaused = !config.TradingPaused

	return config, nil
}

func (r *guildConfigRepo) Update(ctx context.Context, tx pgx.Tx, config *domain.GuildConfig) error {
	const query = `
		UPDATE guild_config
		SET starting_cash = $2,
		    base_price = $3,
		    slope = $4,
		    trade_fee_bps = $5,
		    subject_fee_bps = $6,
		    sunset_days = $7,
		    trading_enabled = $8
		WHERE guild_id = $1
	`

	tradingEnabled := !config.TradingPaused

	_, err := tx.Exec(ctx, query,
		config.GuildID,
		config.StartingCash,
		config.BasePrice,
		config.Slope,
		config.TradeFeeBps,
		config.SubjectFeeBps,
		config.SunsetDays,
		tradingEnabled,
	)

	if err != nil {
		return fmt.Errorf("failed to update guild config: %w", err)
	}

	return nil
}
