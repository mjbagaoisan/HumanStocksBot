package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/mjbagaoisan/humanstocksbot/internal/curve"
	"github.com/mjbagaoisan/humanstocksbot/internal/domain"
	"github.com/mjbagaoisan/humanstocksbot/internal/repos"
)

type QuoteService struct {
	markets repos.MarketRepo
	configs repos.GuildConfigRepo
}

type QuoteResult struct {
	Side        string
	Shares      int64
	GrossAmount int64
	TotalFee    int64
	SubjectFee  int64
	SystemFee   int64
	NetAmount   int64 // cost for buy, payout for sell
	PriceBefore int64
	PriceAfter  int64
}

func NewQuoteService(markets repos.MarketRepo, configs repos.GuildConfigRepo) *QuoteService {
	return &QuoteService{
		markets: markets,
		configs: configs,
	}
}

func (s *QuoteService) Quote(ctx context.Context, tx pgx.Tx, guildID, subjectUserID, side string, qty int64) (*QuoteResult, error) {
	if side != "BUY" && side != "SELL" {
		return nil, fmt.Errorf("invalid side %q: must be BUY or SELL", side)
	}

	if qty <= 0 {
		return nil, domain.ErrInvalidQuantity
	}

	market, err := s.markets.Get(ctx, tx, guildID, subjectUserID)
	if err != nil {
		return nil, err
	}
	if market == nil {
		return nil, domain.ErrMarketNotFound
	}

	config, err := s.configs.GetOrCreate(ctx, tx, guildID)
	if err != nil {
		return nil, err
	}

	if market.Status == domain.MarketStatusClosed {
		return nil, domain.ErrMarketNotActive
	}

	if side == "BUY" && market.Status == domain.MarketStatusSunsetting {
		return nil, domain.ErrMarketSunsetting
	}

	supply := market.SharesOutstanding
	priceBefore := curve.Price(config.BasePrice, config.Slope, supply)

	var grossAmount int64
	var priceAfter int64

	// cakculate gross amount and price after
	if side == "BUY" {
		grossAmount, err = curve.BuyCost(config.BasePrice, config.Slope, supply, qty)
		if err != nil {
			return nil, err
		}
		priceAfter = curve.Price(config.BasePrice, config.Slope, supply+qty)
	} else {
		if qty > supply {
			return nil, domain.ErrInsufficientSupply
		}
		grossAmount, err = curve.SellPayout(config.BasePrice, config.Slope, supply, qty)
		if err != nil {
			return nil, err
		}
		priceAfter = curve.Price(config.BasePrice, config.Slope, supply-qty)
	}

	fees := domain.CalcFees(grossAmount, config.TradeFeeBps, config.SubjectFeeBps)

	var netAmount int64
	if side == "BUY" {
		netAmount = grossAmount + fees.Total
	} else {
		netAmount = grossAmount - fees.Total
	}

	return &QuoteResult{
		Side:        side,
		Shares:      qty,
		GrossAmount: grossAmount,
		TotalFee:    fees.Total,
		SubjectFee:  fees.SubjectFee,
		SystemFee:   fees.SystemFee,
		NetAmount:   netAmount,
		PriceBefore: priceBefore,
		PriceAfter:  priceAfter,
	}, nil
}
