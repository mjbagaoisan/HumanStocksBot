package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/mjbagaoisan/humanstocksbot/internal/config"
	"github.com/mjbagaoisan/humanstocksbot/internal/db"
)

// 1. Load config
// 2. Connect to DB
// 3. Create discord session
// 4. Register slash commands
// 5. Add interaction handler
// 6. Open connection
// 7. Wait for shutdown signal
// 8. Cleanup (close session, close DB pool)

func main() {

	// load config file
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config")
	}

	// connect db
	dbpool, err := db.NewPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	// create discord session
	discord, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Fatalf("failed to create Discord session: %v", err)
	}

	// add interaction handler
	discord.AddHandler(handleInteraction)

	// open connection
	err = discord.Open()
	if err != nil {
		log.Fatalf("failed to open Discord connection: %v", err)
	}

	log.Println("Registering slash commands...")
	registeredCommands, err := discord.ApplicationCommandBulkOverwrite(cfg.DiscordAppID, "", commands)
	if err != nil {
		log.Fatalf("Failed to register commands: %v", err)
	}
	log.Printf("Successfully registered %d commands", len(registeredCommands))

	// wait for shutdown signal
	log.Println("Bot is running. Press CTRL+C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// cleanup
	log.Println("Shutting down...")
	discord.Close()
	dbpool.Close()
}
