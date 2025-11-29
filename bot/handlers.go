package bot

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"amdecrypt/pkg/decrypt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handler manages Telegram bot command handlers.
type Handler struct {
	bot    *tgbotapi.BotAPI
	config *Config

	// Rate limiting
	mu          sync.Mutex
	userLastReq map[int64]time.Time
	rateLimit   time.Duration
}

// NewHandler creates a new Handler instance.
func NewHandler(bot *tgbotapi.BotAPI, config *Config) *Handler {
	return &Handler{
		bot:         bot,
		config:      config,
		userLastReq: make(map[int64]time.Time),
		rateLimit:   30 * time.Second,
	}
}

// checkRateLimit returns true if the user is rate limited.
func (h *Handler) checkRateLimit(userID int64) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	lastReq, exists := h.userLastReq[userID]
	if exists && time.Since(lastReq) < h.rateLimit {
		return true
	}
	h.userLastReq[userID] = time.Now()
	return false
}

// HandleStart handles the /start command.
func (h *Handler) HandleStart(message *tgbotapi.Message) {
	welcomeText := `🎵 *Welcome to Apple Music Decrypt Bot!*

This bot helps you decrypt Apple Music tracks.

*Available Commands:*
• /help - Show detailed usage instructions
• /decrypt <track_id> <key> - Decrypt a track (with file attachment)
• /status - Check bot and server status

To decrypt a track, send the encrypted file along with the /decrypt command.`

	msg := tgbotapi.NewMessage(message.Chat.ID, welcomeText)
	msg.ParseMode = "Markdown"
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("Error sending start message: %v", err)
	}
}

// HandleHelp handles the /help command.
func (h *Handler) HandleHelp(message *tgbotapi.Message) {
	helpText := `📖 *How to use this bot:*

*1. Decrypting a Track:*
Send an encrypted Apple Music file as a document with the caption:
` + "`/decrypt <track_id> <key>`" + `

*Parameters:*
• ` + "`track_id`" + ` - The Apple Music track ID
• ` + "`key`" + ` - The FairPlay Streaming Key URI

*Example:*
Upload your encrypted .m4a file and set the caption to:
` + "`/decrypt 1234567890 skd://itunes.apple.com/...`" + `

*2. Check Status:*
Use /status to verify the bot is working and the server is reachable.

*Notes:*
• Maximum file size: 50MB (Telegram limit)
• Rate limit: 1 request per 30 seconds
• Supported format: Encrypted M4A files`

	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	msg.ParseMode = "Markdown"
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("Error sending help message: %v", err)
	}
}

// HandleStatus handles the /status command.
func (h *Handler) HandleStatus(message *tgbotapi.Message) {
	statusText := fmt.Sprintf(`🟢 *Bot Status: Online*

*Configuration:*
• Agent IP: %s
• MP4Decrypt: %s

Bot is ready to process requests.`, h.config.AgentIP, h.config.MP4DecryptPath)

	msg := tgbotapi.NewMessage(message.Chat.ID, statusText)
	msg.ParseMode = "Markdown"
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("Error sending status message: %v", err)
	}
}

// HandleDecrypt handles the /decrypt command with file upload.
func (h *Handler) HandleDecrypt(message *tgbotapi.Message) {
	userID := message.From.ID

	// Check rate limit
	if h.checkRateLimit(userID) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "⏳ Please wait 30 seconds between requests.")
		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("Error sending rate limit message: %v", err)
		}
		return
	}

	// Parse command arguments
	args := strings.Fields(message.CommandArguments())
	if len(args) < 2 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Usage: /decrypt <track_id> <key>\n\nPlease attach the encrypted file and include both track ID and key.")
		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("Error sending usage message: %v", err)
		}
		return
	}

	trackID := args[0]
	key := args[1]

	// Check for attached document
	if message.Document == nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Please attach the encrypted file as a document with the /decrypt command.")
		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("Error sending no document message: %v", err)
		}
		return
	}

	// Send processing message
	processingMsg := tgbotapi.NewMessage(message.Chat.ID, "⏳ Processing your file...")
	sentMsg, err := h.bot.Send(processingMsg)
	if err != nil {
		log.Printf("Error sending processing message: %v", err)
		return
	}

	// Download the file
	fileConfig := tgbotapi.FileConfig{FileID: message.Document.FileID}
	file, err := h.bot.GetFile(fileConfig)
	if err != nil {
		h.sendError(message.Chat.ID, sentMsg.MessageID, fmt.Sprintf("Failed to get file info: %v", err))
		return
	}

	// Create temporary directory for this request
	tempDir, err := os.MkdirTemp("", "amdecrypt-*")
	if err != nil {
		h.sendError(message.Chat.ID, sentMsg.MessageID, fmt.Sprintf("Failed to create temp directory: %v", err))
		return
	}
	defer os.RemoveAll(tempDir)

	inputPath := filepath.Join(tempDir, "input.m4a")
	outputPath := filepath.Join(tempDir, "output.m4a")

	// Download file content
	fileURL := file.Link(h.bot.Token)
	if err := downloadFile(fileURL, inputPath); err != nil {
		h.sendError(message.Chat.ID, sentMsg.MessageID, fmt.Sprintf("Failed to download file: %v", err))
		return
	}

	// Extract song info
	info, err := decrypt.ExtractSong(inputPath)
	if err != nil {
		h.sendError(message.Chat.ID, sentMsg.MessageID, fmt.Sprintf("Failed to extract song info: %v", err))
		return
	}

	// Decrypt the song
	keys := []string{decrypt.PrefetchKey, key}
	if err := decrypt.DecryptSong(h.config.AgentIP, h.config.MP4DecryptPath, outputPath, trackID, info, keys); err != nil {
		h.sendError(message.Chat.ID, sentMsg.MessageID, fmt.Sprintf("Failed to decrypt: %v", err))
		return
	}

	// Send decrypted file back
	outputFile, err := os.Open(outputPath)
	if err != nil {
		h.sendError(message.Chat.ID, sentMsg.MessageID, fmt.Sprintf("Failed to open output file: %v", err))
		return
	}
	defer outputFile.Close()

	doc := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FileReader{
		Name:   fmt.Sprintf("%s_decrypted.m4a", trackID),
		Reader: outputFile,
	})
	doc.Caption = "✅ Decryption complete!"

	if _, err := h.bot.Send(doc); err != nil {
		h.sendError(message.Chat.ID, sentMsg.MessageID, fmt.Sprintf("Failed to send decrypted file: %v", err))
		return
	}

	// Delete processing message
	deleteMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, sentMsg.MessageID)
	if _, err := h.bot.Request(deleteMsg); err != nil {
		log.Printf("Error deleting processing message: %v", err)
	}
}

// sendError sends an error message and updates the processing message.
func (h *Handler) sendError(chatID int64, messageID int, errorText string) {
	log.Printf("Error: %s", errorText)

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, fmt.Sprintf("❌ Error: %s", errorText))
	if _, err := h.bot.Send(editMsg); err != nil {
		log.Printf("Error sending error message: %v", err)
	}
}
