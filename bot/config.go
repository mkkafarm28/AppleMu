package bot

import (
	"errors"
	"os"
)

// Config holds the bot configuration loaded from environment variables.
type Config struct {
	TelegramBotToken string
	AgentIP          string
	MP4DecryptPath   string
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() (*Config, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return nil, errors.New("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	agentIP := os.Getenv("AGENT_IP")
	if agentIP == "" {
		return nil, errors.New("AGENT_IP environment variable is required")
	}

	mp4DecryptPath := os.Getenv("MP4DECRYPT_PATH")
	if mp4DecryptPath == "" {
		mp4DecryptPath = "mp4decrypt"
	}

	return &Config{
		TelegramBotToken: token,
		AgentIP:          agentIP,
		MP4DecryptPath:   mp4DecryptPath,
	}, nil
}
