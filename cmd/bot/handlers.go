package main

import (
	"github.com/bwmarrin/discordgo"
)

func handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "optin":
		handleOptin(s, i)
	case "optout":
		handleOptout(s, i)
	case "quote":
		handleQuote(s, i)
	case "buy":
		handleBuy(s, i)
	case "sell":
		handleSell(s, i)
	case "portfolio":
		handlePortfolio(s, i)
	case "profile":
		handleProfile(s, i)
	case "leaderboard":
		handleLeaderboard(s, i)
	case "config":
		handleConfig(s, i)
	case "pause":
		handlePause(s, i)
	case "delete_my_data":
		handleDeleteMyData(s, i)
	}
}

// User commands

func handleOptin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	respondPublic(s, i, "Optin command - not implemented yet")
}

func handleOptout(s *discordgo.Session, i *discordgo.InteractionCreate) {
	respondPublic(s, i, "Optout command - not implemented yet")
}

func handleQuote(s *discordgo.Session, i *discordgo.InteractionCreate) {
	respondEphemeral(s, i, "Quote command - not implemented yet")
}

func handleBuy(s *discordgo.Session, i *discordgo.InteractionCreate) {
	respondPublic(s, i, "Buy command - not implemented yet")
}

func handleSell(s *discordgo.Session, i *discordgo.InteractionCreate) {
	respondPublic(s, i, "Sell command - not implemented yet")
}

func handlePortfolio(s *discordgo.Session, i *discordgo.InteractionCreate) {
	respondEphemeral(s, i, "Portfolio command - not implemented yet")
}

func handleProfile(s *discordgo.Session, i *discordgo.InteractionCreate) {
	respondPublic(s, i, "Profile command - not implemented yet")
}

func handleLeaderboard(s *discordgo.Session, i *discordgo.InteractionCreate) {
	respondPublic(s, i, "Leaderboard command - not implemented yet")
}

// Admin commands

func handleConfig(s *discordgo.Session, i *discordgo.InteractionCreate) {
	respondEphemeral(s, i, "Config command - not implemented yet")
}

func handlePause(s *discordgo.Session, i *discordgo.InteractionCreate) {
	respondPublic(s, i, "Pause command - not implemented yet")
}

// Privacy

func handleDeleteMyData(s *discordgo.Session, i *discordgo.InteractionCreate) {
	respondEphemeral(s, i, "Delete my data command - not implemented yet")
}

// Response helpers

func respondPublic(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

func respondEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
