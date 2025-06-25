package main

import (
	"fmt"
	"log"

	"github.com/cordilleradev/hyperliquid-go/pkg/client/stream"
)

// WatchWalletOrders écoute les ordres d'un wallet Hyperliquid et affiche les nouveaux ordres détectés.
func WatchWalletOrders(wallet string) error {
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

	for orders := range wsClient.OrderChan {
		fmt.Printf("Nouvel ordre détecté pour %s: %#v\n", wallet, orders)
	}

	return nil
}

// Pour test rapide : lancer ce fichier seul avec go run watch.go
func mainWatch() {
	wallet := "0xTON_WALLET" // Remplace par l'adresse à surveiller
	if err := WatchWalletOrders(wallet); err != nil {
		log.Fatal(err)
	}
}
