package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjbagaoisan/humanstocksbot/internal/config"
	"github.com/mjbagaoisan/humanstocksbot/internal/domain"
	"github.com/mjbagaoisan/humanstocksbot/internal/repos"
	"github.com/mjbagaoisan/humanstocksbot/internal/services"
)

type Bot struct {
	DB     *pgxpool.Pool
	Config *config.Config
}

func (b *Bot) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in handler: %v\n%s", r, debug.Stack())
			_ = respondEphemeral(s, i, "An internal error occurred")
		}
	}()

	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "optin":
		b.handleOptin(s, i)
	case "optout":
		b.handleOptout(s, i)
	case "quote":
		b.handleQuote(s, i)
	case "buy":
		b.handleBuy(s, i)
	case "sell":
		b.handleSell(s, i)
	case "portfolio":
		b.handlePortfolio(s, i)
	case "profile":
		b.handleProfile(s, i)
	case "leaderboard":
		b.handleLeaderboard(s, i)
	case "config":
		b.handleConfig(s, i)
	case "pause":
		b.handlePause(s, i)
	case "delete_my_data":
		b.handleDeleteMyData(s, i)
	}
}

// User commands

func (b *Bot) handleOptin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()
	guildID := i.GuildID
	userID := ""
	if i.Member != nil && i.Member.User != nil {
		userID = i.Member.User.ID
	} else if i.User != nil {
		userID = i.User.ID
	}

	if guildID == "" || userID == "" {
		if err := respondEphemeral(s, i, "This command must be used in a server."); err != nil {
			log.Printf("failed to respond to optin: %v", err)
		}
		return
	}

	service := services.NewOptInService(
		b.DB,
		repos.NewGuildConfigRepo(b.DB),
		repos.NewGuildMemberRepo(b.DB),
		repos.NewWalletRepo(b.DB),
		repos.NewMarketRepo(b.DB),
	)

	result, err := service.OptIn(ctx, guildID, userID)
	if err != nil {
		var message string
		if errors.Is(err, domain.ErrAlreadyOptedIn) {
			message = "You have already opted in."
		} else {
			message = fmt.Sprintf("Failed to opt in: %v", err)
		}

		if respondErr := respondEphemeral(s, i, message); respondErr != nil {
			log.Printf("failed to respond to optin error: %v", respondErr)
		}
		return
	}

	content := fmt.Sprintf("**<@%s>** is now tradable! Starting price: **%s**", result.UserID, formatDollars(result.BasePrice))
	if err := respondPublic(s, i, content); err != nil {
		log.Printf("failed to respond to optin: %v", err)
	}
}

func (b *Bot) handleOptout(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := respondPublic(s, i, "Optout command - not implemented yet"); err != nil {
		log.Printf("failed to respond to optout: %v", err)
	}
}

func (b *Bot) handleQuote(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()
	guildID := i.GuildID
	opts := i.ApplicationCommandData().Options

	var subjectUserID string
	var subjectName string
	var qty int64
	var side string

	resolved := i.ApplicationCommandData().Resolved
	for _, opt := range opts {
		switch opt.Name {
		case "user":
			subjectUserID = opt.UserValue(nil).ID
			if resolved != nil && resolved.Users != nil {
				if u, ok := resolved.Users[subjectUserID]; ok {
					subjectName = u.GlobalName
					if subjectName == "" {
						subjectName = u.Username
					}
				}
			}
		case "quantity":
			qty = opt.IntValue()
		case "action":
			side = strings.ToUpper(opt.StringValue())
		}
	}

	tx, err := b.DB.Begin(ctx)
	if err != nil {
		log.Printf("failed to begin tx for quote: %v", err)
		_ = respondEphemeral(s, i, "Something went wrong. Please try again.")
		return
	}
	defer tx.Rollback(ctx)

	svc := services.NewQuoteService(
		repos.NewMarketRepo(b.DB),
		repos.NewGuildConfigRepo(b.DB),
	)

	result, err := svc.Quote(ctx, tx, guildID, subjectUserID, side, qty)
	if err != nil {
		var message string
		switch {
		case errors.Is(err, domain.ErrInvalidQuantity):
			message = "Quantity must be greater than zero."
		case errors.Is(err, domain.ErrMarketNotFound):
			message = "That user hasn't opted in yet."
		case errors.Is(err, domain.ErrMarketNotActive):
			message = "That market is closed."
		case errors.Is(err, domain.ErrMarketSunsetting):
			message = "That market is sunsetting — only sells are allowed."
		case errors.Is(err, domain.ErrInsufficientSupply):
			message = "Not enough shares outstanding to sell that many."
		default:
			message = fmt.Sprintf("Failed to get quote: %v", err)
		}
		_ = respondEphemeral(s, i, message)
		return
	}

	totalLabel := "Total cost"
	color := 0x2ECC71 // green for buy
	if result.Side == "SELL" {
		totalLabel = "Payout"
		color = 0xE74C3C // red for sell
	}

	feePercent := float64(result.TotalFee) / float64(result.GrossAmount) * 100

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Quote: %s %d shares of %s", result.Side, result.Shares, subjectName),
		Color: color,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Gross", Value: formatDollars(result.GrossAmount), Inline: true},
			{Name: fmt.Sprintf("Fee (%.1f%%)", feePercent), Value: formatDollars(result.TotalFee), Inline: true},
			{Name: totalLabel, Value: fmt.Sprintf("**%s**", formatDollars(result.NetAmount)), Inline: true},
			{Name: "Price before", Value: formatDollars(result.PriceBefore), Inline: true},
			{Name: "Price after", Value: formatDollars(result.PriceAfter), Inline: true},
			{Name: "Fee split", Value: fmt.Sprintf("Subject: %s | System: %s", formatDollars(result.SubjectFee), formatDollars(result.SystemFee)), Inline: false},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Prices are estimates and may change before execution",
		},
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		log.Printf("failed to respond to quote: %v", err)
	}
}

func (b *Bot) handleBuy(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := respondPublic(s, i, "Buy command - not implemented yet"); err != nil {
		log.Printf("failed to respond to buy: %v", err)
	}
}

func (b *Bot) handleSell(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := respondPublic(s, i, "Sell command - not implemented yet"); err != nil {
		log.Printf("failed to respond to sell: %v", err)
	}
}

func (b *Bot) handlePortfolio(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := respondEphemeral(s, i, "Portfolio command - not implemented yet"); err != nil {
		log.Printf("failed to respond to portfolio: %v", err)
	}
}

func (b *Bot) handleProfile(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := respondPublic(s, i, "Profile command - not implemented yet"); err != nil {
		log.Printf("failed to respond to profile: %v", err)
	}
}

func (b *Bot) handleLeaderboard(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := respondPublic(s, i, "Leaderboard command - not implemented yet"); err != nil {
		log.Printf("failed to respond to leaderboard: %v", err)
	}
}

// Admin commands

func (b *Bot) handleConfig(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := respondEphemeral(s, i, "Config command - not implemented yet"); err != nil {
		log.Printf("failed to respond to config: %v", err)
	}
}

func (b *Bot) handlePause(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := respondPublic(s, i, "Pause command - not implemented yet"); err != nil {
		log.Printf("failed to respond to pause: %v", err)
	}
}

// Privacy

func (b *Bot) handleDeleteMyData(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := respondEphemeral(s, i, "Delete my data command - not implemented yet"); err != nil {
		log.Printf("failed to respond to delete_my_data: %v", err)
	}
}

// Response helpers

func respondPublic(s *discordgo.Session, i *discordgo.InteractionCreate, content string) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

func respondEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, content string) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func formatDollars(cents int64) string {
	dollars := cents / 100
	remainder := cents % 100
	if remainder < 0 {
		remainder = -remainder
	}
	return fmt.Sprintf("$%d.%02d", dollars, remainder)
}
