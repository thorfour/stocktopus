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
