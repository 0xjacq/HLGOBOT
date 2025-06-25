# HLGOBOT

Go bot for Hyperliquid:
- **Wallet monitoring** with real-time Telegram notifications
- **Leaderboard Finder**: find the traders closest to target criteria

---

## 🚀 Installation

1. **Clone the repo**
   ```sh
   git clone https://github.com/youruser/HLGOBOT.git
   cd HLGOBOT
   ```
2. **Install dependencies**
   ```sh
   go mod tidy
   ```

---

## ⚙️ Configuration

Create a `config.json` file at the root:
```json
{
  "bot_token": "YOUR_BOT_TOKEN",
  "chat_id": "YOUR_CHAT_ID"
}
```
- Get a token via [@BotFather](https://t.me/BotFather) on Telegram
- Retrieve your chat_id via the Telegram API

---

## 🕵️‍♂️ Monitor a Hyperliquid wallet (watch)

Run the watcher to receive a Telegram notification for every new order:
```sh
# From the project root
# Replace <wallet_address> with the address you want to monitor

go run cmd/watch/main.go <wallet_address>
```

---

## 🏆 Use the Leaderboard Finder

Find the N traders closest to your target criteria (equity, pnl, volume, etc.):
```sh
# Example usage with custom flags

go run cmd/leaderboardfinder/main.go \
  -equity=138526.48 \
  -pnl=173310.97 \
  -volume=11801800.46 \
  -period=alltime \
  -top=2
```

---

## 🗂️ Project Structure

```
HLGOBOT/
│
├── cmd/
│   ├── watch/                # CLI for the Telegram watcher
│   │   └── main.go
│   └── leaderboardfinder/    # CLI for the leaderboard finder
│       └── main.go
│
├── internal/
│   ├── watcher/              # Watcher business logic
│   │   └── watch.go
│   └── leaderboard/          # Leaderboard finder business logic
│       └── finder.go
│
├── config.json               # Telegram configuration
├── go.mod / go.sum           # Go dependencies
└── .gitignore
```

---

## 👩‍💻 Contributing

- Fork the repo, create a branch, submit a PR!
- Add new CLI modules in `cmd/` for other use cases (stats, export, etc.)
- Add unit tests in `internal/`
- Suggest improvements (multi-wallet, multi-chat, other channels...)

---

## 📚 Resources
- [Hyperliquid API Go](https://pkg.go.dev/github.com/cordilleradev/hyperliquid-go)
- [Telegram Bot API](https://core.telegram.org/bots/api)

---

## 📝 License
MIT 