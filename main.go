package main

import (
	"HLGOBOT/leaderboard"
	"fmt"
	"os"
)

func main() {
	// Example: Simulate leaderboard data (replace with real data fetching in production)
	traders := []leaderboard.TraderPerformance{
		{
			EthAddress:   "0xabc",
			AccountValue: 100000,
			AllTime:      leaderboard.PeriodPerformance{Pnl: 20000, Vlm: 500000, Roi: 0.2},
		},
		{
			EthAddress:   "0xdef",
			AccountValue: 138526.48,
			AllTime:      leaderboard.PeriodPerformance{Pnl: 173310.97, Vlm: 11801800.46, Roi: 0.3},
		},
		{
			EthAddress:   "0x123",
			AccountValue: 140000,
			AllTime:      leaderboard.PeriodPerformance{Pnl: 170000, Vlm: 12000000, Roi: 0.25},
		},
	}

	targetEquity := 138526.48
	targetPnL := 173310.97
	targetVolume := 11801800.46
	period := "alltime"
	topN := 2

	closest := leaderboard.FindClosestTraders(traders, targetEquity, targetPnL, targetVolume, period, topN)
	fmt.Printf("Top %d closest addresses to the target values for period '%s':\n", topN, period)
	for i, t := range closest {
		p := leaderboard.GetPeriod(&t, period)
		fmt.Printf("%d. Address: %s\n", i+1, t.EthAddress)
		fmt.Printf("   Equity = %.2f, PnL = %.2f, Volume = %.2f\n", t.AccountValue, p.Pnl, p.Vlm)
	}

	os.Exit(0)
}
