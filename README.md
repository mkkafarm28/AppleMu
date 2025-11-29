# amdecrypt

A CLI tool for decrypting Apple Music songs in conjunction with a [wrapper server](https://github.com/WorldObservationLog/wrapper).

Based on [apple-music-alac-downloader](https://github.com/alacleaker/apple-music-alac-downloader).

## 📋 Prerequisites

Must be added to your system's PATH:

- [mp4decrypt](https://www.bento4.com/downloads/)

## 📦 Installation

1. Download the latest binary for your platform from the [releases page](https://github.com/glomatico/amdecrypt/releases/latest)
2. Extract the archive
3. Add the binary to your system's PATH

## 🚀 Usage

This tool is designed to be called by other programs rather than used directly. For example, [gamdl](https://github.com/glomatico/gamdl) will automatically invoke amdecrypt when needed.

### Manual Usage

If needed, you can also run it directly from the command line:

```bash
amdecrypt <agentIp> <mp4decryptPath> <id> <key> <inputPath> <outputPath>
```

### Arguments

| Argument         | Description                      |
| ---------------- | -------------------------------- |
| `agentIp`        | IP address of the wrapper server |
| `mp4decryptPath` | Path to the mp4decrypt binary    |
| `id`             | Track ID                         |
| `key`            | FairPlay Streaming Key           |
| `inputPath`      | Path to the encrypted file       |
| `outputPath`     | Path for the decrypted output    |

## 🤖 Telegram Bot

The project includes a Telegram bot that allows users to decrypt Apple Music tracks through Telegram.

### Bot Setup

1. Create a new bot with [@BotFather](https://t.me/BotFather) on Telegram
2. Copy the bot token provided by BotFather
3. Configure environment variables (see below)
4. Run the bot

### Environment Variables

Create a `.env` file or set the following environment variables:

| Variable             | Required | Description                                  | Default      |
| -------------------- | -------- | -------------------------------------------- | ------------ |
| `TELEGRAM_BOT_TOKEN` | Yes      | Bot token from BotFather                     | -            |
| `AGENT_IP`           | Yes      | IP address and port of the wrapper server    | -            |
| `MP4DECRYPT_PATH`    | No       | Path to mp4decrypt binary                    | `mp4decrypt` |

Example `.env` file:
```bash
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
AGENT_IP=127.0.0.1:10020
MP4DECRYPT_PATH=mp4decrypt
```

### Running the Bot

```bash
# Build the bot
go build -o bot ./cmd/bot

# Run with environment variables
TELEGRAM_BOT_TOKEN=your_token AGENT_IP=127.0.0.1:10020 ./bot
```

### Bot Commands

| Command                      | Description                                    |
| ---------------------------- | ---------------------------------------------- |
| `/start`                     | Welcome message and usage instructions         |
| `/help`                      | Show available commands and how to use them    |
| `/decrypt <track_id> <key>`  | Decrypt a track (attach file with this caption)|
| `/status`                    | Check bot and server status                    |

### How to Decrypt a Track

1. Open the bot in Telegram
2. Send the encrypted `.m4a` file as a document
3. Use `/decrypt <track_id> <key>` as the caption
4. The bot will process the file and send back the decrypted version

**Note:** Rate limiting is enabled (1 request per 30 seconds per user).

## 🏗️ Project Structure

```
.
├── main.go              # CLI entry point
├── cmd/
│   └── bot/
│       └── main.go      # Telegram bot entry point
├── bot/
│   ├── bot.go           # Main bot logic
│   ├── config.go        # Configuration management
│   └── handlers.go      # Command handlers
├── pkg/
│   └── decrypt/
│       └── decrypt.go   # Core decryption functions
├── .env.example         # Example environment configuration
└── README.md
```

## ⚠️ Disclaimer

This tool was mostly created with AI assistance.
