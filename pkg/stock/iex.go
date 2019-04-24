package stock

import (
	iex "github.com/thorfour/iex/pkg/api"
	iextype "github.com/thorfour/iex/pkg/types"
)

// IexWrapper is a wrapper around the IEX library
type IexWrapper struct{}

// Price returns the current price of the ticker
func (w *IexWrapper) Price(ticker string) (float64, error) {
	return iex.Price(ticker)
}

// BatchQuotes returns a slice of quotes for the given tickers
func (w *IexWrapper) BatchQuotes(tickers []string) ([]*Quote, error) {
	batch, err := iex.BatchQuotes(tickers)
	if err != nil {
		return nil, err
	}

	var quotes []*Quote
	for ticker := range batch {
		q, err := batch.Quote(ticker)
		if err != nil {
			return nil, err
		}

		quotes = append(quotes, &Quote{
			Ticker:        ticker,
			LatestPrice:   q.LatestPrice,
			Change:        q.Change,
			ChangePercent: q.ChangePercent,
		})
	}

	return quotes, nil
}

// News returns recent news for a ticker
func (w *IexWrapper) News(ticker string) ([]string, error) {
	latest, err := iex.News(ticker)
	if err != nil {
		return nil, err
	}

	var news []string
	for _, n := range latest {
		news = append(news, n.Summary)
	}

	return news, nil
}

// Stats returns the stats for a ticker
func (w *IexWrapper) Stats(ticker string) (*iextype.Stats, error) {
	stats, err := iex.Stats(ticker)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
