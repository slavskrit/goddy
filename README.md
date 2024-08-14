# 🤖 Goddy Bot

Goddy Bot is a simple Go-based Telegram bot that can download and send videos from Instagram and X (formerly Twitter) links shared in a chat. It also reacts to messages with an emoji and can process voice messages.

## 🚀 Features

- 🔗 Detects and processes Instagram and X links in chat messages
- 🎥 Downloads and sends videos directly in the chat
- 🤔 Reacts to messages with a thinking emoji

## 🛠️ Installation

### Clone the Repository

```bash
git clone https://github.com/yourusername/goddy-bot.git
cd goddy-bot
```

### Build the Bot
```bash
docker build -t goddy-bot .
```

### Run the Bot
```bash
docker run -d --name goddy-bot -e TELEGRAM_TOKEN=your_telegram_bot_token goddy-bot
```

## 🐳 Docker
Goddy Bot is containerized using Docker. The Dockerfile builds the Go binary, installs yt-dlp in a Python virtual environment, and runs the bot.

### Docker Commands
Build the image:

```bash
docker build -t goddy-bot .
```

### Run the container:

```bash
docker run -d --name goddy-bot -e TELEGRAM_TOKEN=your_telegram_bot_token goddy-bot
```

### Stop the container:

```bash
docker stop goddy-bot
```

### Remove the container:

```bash
docker rm goddy-bot
```

### 📄 Environment Variables
TELEGRAM_TOKEN: Your Telegram bot token obtained from BotFather.

## 🔧 Configuration
You can modify the bot to handle additional types of messages or integrate with other platforms by editing the MainHandler and processMessage functions in the code.

## 📝 License
This project is licensed under the MIT License - see the LICENSE file for details.

## 👨‍💻 Contributing
Feel free to submit issues or pull requests. Contributions are always welcome! 🎉
