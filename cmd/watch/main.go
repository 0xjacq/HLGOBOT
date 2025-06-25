package main

import (
	"fmt"
	"log"
	"os"

	"github.com/0xjacq/HLGOBOT/internal/watcher"
)

func main() {
	cfg, err := watcher.LoadTelegramConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// For now, only monitoring a single wallet (can be improved for multiple wallets)
	if len(os.Args) < 2 {
		fmt.Println("Usage: watch <wallet_address>")
		os.Exit(1)
	}
	wallet := os.Args[1]

	if err := watcher.WatchWalletOrders(cfg, wallet); err != nil {
		log.Fatal(err)
	}
}
