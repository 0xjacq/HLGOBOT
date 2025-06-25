package watcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/cordilleradev/hyperliquid-go/pkg/client/stream"
)

type TelegramConfig struct {
	BotToken string `json:"bot_token"`
	ChatID   string `json:"chat_id"`
}

func LoadTelegramConfig(path string) (*TelegramConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg TelegramConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SendTelegramMessage(token, chatID, message string) error {
	apiURL := "https://api.telegram.org/bot" + token + "/sendMessage"
	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id": {chatID},
		"text":    {message},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// WatchWalletOrders listens for orders on a Hyperliquid wallet and sends a Telegram notification for each new order.
func WatchWalletOrders(cfg *TelegramConfig, wallet string) error {
	wsClient, err := stream.NewHyperliquidWebsocketClient("wss://api.hyperliquid.xyz/ws")
	if err != nil {
		return fmt.Errorf("error creating websocket client: %w", err)
	}

	go func() {
		for err := range wsClient.ErrorChan {
			log.Println("WebSocket error:", err)
		}
	}()

	err = wsClient.StreamOrderUpdates(wallet)
	if err != nil {
		return fmt.Errorf("error streaming order updates: %w", err)
	}

	fmt.Printf("Watcher started for %s, waiting for orders...\n", wallet)
	for orders := range wsClient.OrderChan {
		for _, order := range orders {
			o := order.Order
			var side string
			switch o.Side {
			case "A":
				side = "BUY"
			case "B":
				side = "SELL"
			default:
				side = o.Side
			}
			msg := fmt.Sprintf("[%s] %s %s %s @ %s (oid: %d, status: %s)",
				wallet, side, o.Sz, o.Coin, o.LimitPx, o.Oid, order.Status)
			fmt.Println(msg) // local debug
			if cfg != nil && cfg.BotToken != "" && cfg.ChatID != "" {
				err := SendTelegramMessage(cfg.BotToken, cfg.ChatID, msg)
				if err != nil {
					log.Println("Telegram send error:", err)
				}
			}
		}
	}

	return nil
}
