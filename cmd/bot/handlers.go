package main

import (
	"log"
	"runtime/debug"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjbagaoisan/humanstocksbot/internal/config"
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
	if err := respondPublic(s, i, "Optin command - not implemented yet"); err != nil {
		log.Printf("failed to respond to optin: %v", err)
	}
}

func (b *Bot) handleOptout(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := respondPublic(s, i, "Optout command - not implemented yet"); err != nil {
		log.Printf("failed to respond to optout: %v", err)
	}
}

func (b *Bot) handleQuote(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := respondEphemeral(s, i, "Quote command - not implemented yet"); err != nil {
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
