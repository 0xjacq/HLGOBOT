// Package leaderboard provides utilities to find the closest traders to target values from Hyperliquid leaderboard data.
package leaderboard

import (
	"math"
	"sort"
	"strings"
)

// PeriodPerformance holds PnL, ROI, and Volume for a period
type PeriodPerformance struct {
	Pnl float64
	Roi float64
	Vlm float64
}

// TraderPerformance holds trader stats
type TraderPerformance struct {
	EthAddress   string
	AccountValue float64
	Day          PeriodPerformance
	Week         PeriodPerformance
	Month        PeriodPerformance
	AllTime      PeriodPerformance
}

// Result holds the index and distance for sorting
type Result struct {
	Idx  int
	Dist float64
}

// GetPeriod returns the PeriodPerformance for the given period string
func GetPeriod(t *TraderPerformance, period string) PeriodPerformance {
	switch strings.ToLower(period) {
	case "day":
		return t.Day
	case "week":
		return t.Week
	case "month", "30d":
		return t.Month
	case "alltime":
		return t.AllTime
	default:
		return t.AllTime
	}
}

// FindClosestTraders finds the top N traders closest to the target values for the given period
func FindClosestTraders(traders []TraderPerformance, targetEquity, targetPnL, targetVolume float64, period string, topN int) []TraderPerformance {
	var results []Result
	for idx, t := range traders {
		p := GetPeriod(&t, period)
		dist := math.Abs(p.Pnl-targetPnL) +
			math.Abs(p.Vlm-targetVolume) +
			math.Abs(t.AccountValue-targetEquity)
		results = append(results, Result{idx, dist})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Dist < results[j].Dist
	})
	var closest []TraderPerformance
	for i := 0; i < topN && i < len(results); i++ {
		closest = append(closest, traders[results[i].Idx])
	}
	return closest
}
