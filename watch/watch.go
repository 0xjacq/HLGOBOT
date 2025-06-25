package watch

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

func loadTelegramConfig(path string) (*TelegramConfig, error) {
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

func sendTelegramMessage(token, chatID, message string) error {
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

// WatchWalletOrders écoute les ordres d'un wallet Hyperliquid et envoie une notif Telegram à chaque nouvel ordre.
func WatchWalletOrders(wallet string) error {
	cfg, err := loadTelegramConfig("config.json")
	if err != nil {
		log.Println("Impossible de charger config.json, notifications Telegram désactivées:", err)
	}

	wsClient, err := stream.NewHyperliquidWebsocketClient("wss://api.hyperliquid.xyz/ws")
	if err != nil {
		return fmt.Errorf("erreur création client websocket: %w", err)
	}

	go func() {
		for err := range wsClient.ErrorChan {
			log.Println("Erreur WebSocket:", err)
		}
	}()

	err = wsClient.StreamOrderUpdates(wallet)
	if err != nil {
		return fmt.Errorf("erreur stream order updates: %w", err)
	}

	fmt.Printf("Watcher démarré pour %s, en attente d'ordres...\n", wallet)
	for orders := range wsClient.OrderChan {
		for _, order := range orders {
			o := order.Order
			var side string
			if o.Side == "A" {
				side = "ACHAT"
			} else if o.Side == "B" {
				side = "VENTE"
			} else {
				side = o.Side
			}
			msg := fmt.Sprintf("[%s] %s %s %s @ %s (oid: %d, status: %s)",
				wallet, side, o.Sz, o.Coin, o.LimitPx, o.Oid, order.Status)
			fmt.Println(msg) // debug local
			if cfg != nil && cfg.BotToken != "" && cfg.ChatID != "" {
				err := sendTelegramMessage(cfg.BotToken, cfg.ChatID, msg)
				if err != nil {
					log.Println("Erreur envoi Telegram:", err)
				}
			}
		}
	}

	return nil
}
