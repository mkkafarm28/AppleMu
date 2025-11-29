package bot

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// httpClient is a shared HTTP client with timeout configuration.
var httpClient = &http.Client{
	Timeout: 5 * time.Minute, // Allow up to 5 minutes for large file downloads
}

// Bot represents the Telegram bot instance.
type Bot struct {
	api     *tgbotapi.BotAPI
	config  *Config
	handler *Handler
}

// New creates a new Bot instance with the given configuration.
func New(config *Config) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	log.Printf("Authorized on account %s", api.Self.UserName)

	return &Bot{
		api:     api,
		config:  config,
		handler: NewHandler(api, config),
	}, nil
}

// Start starts the bot and begins processing updates.
func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	log.Println("Bot started. Listening for updates...")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
		}
	}

	return nil
}

// handleCommand routes commands to the appropriate handler.
func (b *Bot) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		b.handler.HandleStart(message)
	case "help":
		b.handler.HandleHelp(message)
	case "status":
		b.handler.HandleStatus(message)
	case "decrypt":
		b.handler.HandleDecrypt(message)
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Unknown command. Use /help to see available commands.")
		if _, err := b.api.Send(msg); err != nil {
			log.Printf("Error sending unknown command message: %v", err)
		}
	}
}

// downloadFile downloads a file from a URL and saves it to the specified path.
func downloadFile(url, filepath string) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
