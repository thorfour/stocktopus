package stock

// Quote holds all the information for a stock quote
type Quote struct {
	// LatestPrice of the stock
	LatestPrice float64
	// Change in price
	Change float64
	// ChangePercent daily percent change
	ChangePercent float64
}

// Lookup is the interface for a package to do stock lookups
type Lookup interface {
	BatchQuotes([]string) ([]Quote, error)
	Price(string) (int, error)
	News(string) ([]string, error)
}
