package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/0xjacq/HLGOBOT/leaderboard"
	"github.com/cordilleradev/hyperliquid-go/pkg/client/rest"
)

func main() {
	// Définir les flags pour les arguments de la ligne de commande
	targetEquity := flag.Float64("equity", 138526.48, "Target equity value")
	targetPnL := flag.Float64("pnl", 173310.97, "Target PnL value")
	targetVolume := flag.Float64("volume", 11801800.46, "Target volume value")
	period := flag.String("period", "alltime", "Period (day, week, month, alltime)")
	topN := flag.Int("top", 2, "Number of closest traders to show")
	flag.Parse()

	// Création du client Hyperliquid
	client, err := rest.NewHyperliquidRestClient("https://api.hyperliquid.xyz", "https://stats-data.hyperliquid.xyz/Mainnet/leaderboard")
	if err != nil {
		log.Fatalf("Erreur lors de la création du client Hyperliquid: %v", err)
	}

	// Récupération des données du leaderboard
	leaderboardData, err := client.LeaderboardCall()
	if err != nil {
		log.Fatalf("Erreur lors de la récupération du leaderboard: %v", err)
	}

	// Conversion des données du leaderboard en []leaderboard.TraderPerformance
	var traders []leaderboard.TraderPerformance
	for _, entry := range leaderboardData {
		traders = append(traders, leaderboard.TraderPerformance{
			EthAddress:   entry.EthAddress,
			AccountValue: entry.AccountValue,
			Day:          leaderboard.PeriodPerformance(entry.Day),
			Week:         leaderboard.PeriodPerformance(entry.Week),
			Month:        leaderboard.PeriodPerformance(entry.Month),
			AllTime:      leaderboard.PeriodPerformance(entry.AllTime),
		})
	}

	closest := leaderboard.FindClosestTraders(traders, *targetEquity, *targetPnL, *targetVolume, *period, *topN)
	fmt.Printf("Top %d closest addresses to the target values for period '%s':\n", *topN, *period)
	for i, t := range closest {
		p := leaderboard.GetPeriod(&t, *period)
		fmt.Printf("%d. Address: %s\n", i+1, t.EthAddress)
		fmt.Printf("   Equity = %.2f, PnL = %.2f, Volume = %.2f\n", t.AccountValue, p.Pnl, p.Vlm)
	}

	os.Exit(0)
}
