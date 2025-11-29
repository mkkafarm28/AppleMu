package main

import (
	"log"

	"amdecrypt/bot"
)

func main() {
	config, err := bot.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	b, err := bot.New(config)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	if err := b.Start(); err != nil {
		log.Fatalf("Bot stopped with error: %v", err)
	}
}
