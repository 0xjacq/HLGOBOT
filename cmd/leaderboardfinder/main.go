package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/0xjacq/HLGOBOT/internal/leaderboard"
	"github.com/cordilleradev/hyperliquid-go/pkg/client/rest"
)

func main() {
	// Define flags for command-line arguments
	targetEquity := flag.Float64("equity", 138526.48, "Target equity value")
	targetPnL := flag.Float64("pnl", 173310.97, "Target PnL value")
	targetVolume := flag.Float64("volume", 11801800.46, "Target volume value")
	period := flag.String("period", "alltime", "Period (day, week, month, alltime)")
	topN := flag.Int("top", 2, "Number of closest traders to show")
	flag.Parse()

	// Create Hyperliquid client
	client, err := rest.NewHyperliquidRestClient("https://api.hyperliquid.xyz", "https://stats-data.hyperliquid.xyz/Mainnet/leaderboard")
	if err != nil {
		log.Fatalf("Error creating Hyperliquid client: %v", err)
	}

	// Fetch leaderboard data
	leaderboardData, err := client.LeaderboardCall()
	if err != nil {
		log.Fatalf("Error fetching leaderboard: %v", err)
	}

	// Map fields
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
