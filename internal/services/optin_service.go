package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjbagaoisan/humanstocksbot/internal/domain"
	"github.com/mjbagaoisan/humanstocksbot/internal/repos"
)

// OptInService handles user opt-in to the market
type OptInService struct {
	db      *pgxpool.Pool
	configs repos.GuildConfigRepo
	members repos.GuildMemberRepo
	wallets repos.WalletRepo
	markets repos.MarketRepo
}

type OptInResult struct {
	UserID       string
	StartingCash int64
	BasePrice    int64
}

func NewOptInService(
	db *pgxpool.Pool,
	configs repos.GuildConfigRepo,
	members repos.GuildMemberRepo,
	wallets repos.WalletRepo,
	markets repos.MarketRepo,
) *OptInService {
	return &OptInService{
		db:      db,
		configs: configs,
		members: members,
		wallets: wallets,
		markets: markets,
	}
}

func (s *OptInService) OptIn(ctx context.Context, guildID, userID string) (*OptInResult, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	config, err := s.configs.GetOrCreate(ctx, tx, guildID)
	if err != nil {
		return nil, err
	}

	member, err := s.members.Get(ctx, tx, guildID, userID)
	if err != nil {
		return nil, err
	}

	if member != nil && member.OptedIn {
		return nil, domain.ErrAlreadyOptedIn
	}

	if member == nil {
		if err := s.members.Create(ctx, tx, guildID, userID); err != nil {
			return nil, err
		}
	}

	if err := s.members.SetOptedIn(ctx, tx, guildID, userID, true); err != nil {
		return nil, err
	}

	if err := s.wallets.Create(ctx, tx, guildID, userID, config.StartingCash); err != nil {
		return nil, err
	}

	existingMarket, err := s.markets.Get(ctx, tx, guildID, userID)
	if err != nil {
		return nil, err
	}

	if existingMarket == nil {
		market := &domain.Market{
			GuildID:           guildID,
			SubjectUserID:     userID,
			Status:            domain.MarketStatusActive,
			SharesOutstanding: 0,
			ReserveBalance:    0,
			LastPrice:         config.BasePrice,
			SunsetAt:          nil,
		}
		if err := s.markets.Create(ctx, tx, market); err != nil {
			return nil, err
		}
	} else {
		existingMarket.Status = domain.MarketStatusActive
		existingMarket.SharesOutstanding = 0
		existingMarket.ReserveBalance = 0
		existingMarket.LastPrice = config.BasePrice
		existingMarket.SunsetAt = nil
		if err := s.markets.Update(ctx, tx, existingMarket); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	committed = true

	return &OptInResult{
		UserID:       userID,
		StartingCash: config.StartingCash,
		BasePrice:    config.BasePrice,
	}, nil
}
