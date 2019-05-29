package stock

import (
	"github.com/thorfour/iex/pkg/types"
)

// Quote holds all the information for a stock quote
type Quote struct {
	// Ticker is the ticker with the associated quote
	Ticker string
	// LatestPrice of the stock
	LatestPrice float64
	// Change in price
	Change float64
	// ChangePercent daily percent change
	ChangePercent float64
}

// Lookup is the interface for a package to do stock lookups
type Lookup interface {
	BatchQuotes([]string) ([]*Quote, error)
	Price(string) (float64, error)
	News(string) ([]string, error)
	Stats(string) (*types.Stats, error)
	Company(string) (*types.Company, error)
}

// StatsToRows converts a stats struct into a label list of printable values
func StatsToRows(s *types.Stats) [][]interface{} {
	return [][]interface{}{
		{"Marketcap", s.Marketcap},
		{"Beta", s.Beta},
		{"52WeekHigh", s.Week52High},
		{"52WeekLow", s.Week52Low},
		{"DividendRate", s.DividendRate},
		{"LatestEPS", s.LatestEPS},
		{"200 SMA", s.Day200MovingAvg},
		{"50 SMA", s.Day50MovingAvg},
	}
}
