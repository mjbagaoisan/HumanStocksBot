package main

import "github.com/bwmarrin/discordgo"

var commands = []*discordgo.ApplicationCommand{
	// User commands
	{
		Name:        "optin",
		Description: "Opt in to the human stocks market",
	},
	{
		Name:        "optout",
		Description: "Opt out of the human stocks market",
	},
	{
		Name:        "quote",
		Description: "Get a quote for a human stock",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to get a quote for",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "quantity",
				Description: "Number of shares",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "action",
				Description: "Buy or sell",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "buy", Value: "buy"},
					{Name: "sell", Value: "sell"},
				},
			},
		},
	},
	{
		Name:        "buy",
		Description: "Buy shares of a human stock",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to buy shares of",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "quantity",
				Description: "Number of shares to buy",
				Required:    true,
			},
		},
	},
	{
		Name:        "sell",
		Description: "Sell shares of a human stock",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to sell shares of",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "quantity",
				Description: "Number of shares to sell",
				Required:    true,
			},
		},
	},
	{
		Name:        "portfolio",
		Description: "View your portfolio",
	},
	{
		Name:        "profile",
		Description: "View a user's market profile",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to view",
				Required:    true,
			},
		},
	},
	{
		Name:        "leaderboard",
		Description: "View the leaderboard",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "sort",
				Description: "Sort by price or networth",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "price", Value: "price"},
					{Name: "networth", Value: "networth"},
				},
			},
		},
	},

	// Admin commands
	{
		Name:        "config",
		Description: "Configure guild settings (admin only)",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "starting_balance",
				Description: "Starting balance for new users",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "base_price",
				Description: "Base price for new markets",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "slope",
				Description: "Slope for bonding curve",
				Required:    false,
			},
		},
	},
	{
		Name:        "pause",
		Description: "Pause or unpause trading (admin only)",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "paused",
				Description: "True to pause, false to unpause",
				Required:    true,
			},
		},
	},

	// Privacy
	{
		Name:        "delete_my_data",
		Description: "Delete all your data from the bot",
	},
}
