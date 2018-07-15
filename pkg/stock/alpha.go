package stock

import (
	"sync"

	av "github.com/cmckee-dev/go-alpha-vantage"
	iex "github.com/thorfour/iex/pkg/api"
)

// AlphaWrapper is a wrapper around the AlphaVantage library
type AlphaWrapper struct {
	// APIKey is the API key from alpha vantage
	APIKey string
}

// Price reutrns the current price of the ticker
func (w *AlphaWrapper) Price(ticker string) (float64, error) {
	client := av.NewClient(w.APIKey)

	series, err := client.StockTimeSeriesIntraday(av.TimeIntervalOneMinute, ticker)
	if err != nil {
		return 0, err
	}

	return series[0].Close, nil
}

// BatchQuotes returns a slice of quotes for the given tickers
func (w *AlphaWrapper) BatchQuotes(tickers []string) ([]*Quote, error) {
	client := av.NewClient(w.APIKey)

	// AlphaVantage doesn't provide batch requests, make them all in parallel
	resp := make(chan *Quote, len(tickers))
	errCh := make(chan error, len(tickers))
	wg := new(sync.WaitGroup)
	wg.Add(len(tickers))
	for _, ticker := range tickers {
		go func(symbol string) {
			defer wg.Done()
			series, err := client.StockTimeSeriesIntraday(av.TimeIntervalOneMinute, symbol)
			if err != nil {
				errCh <- err
				return
			}

			// convert the series into a quote object
			q := &Quote{
				Ticker:        symbol,
				LatestPrice:   series[0].Close,
				Change:        series[len(series)-1].Close - series[0].Close,
				ChangePercent: ((series[0].Close - series[len(series)-1].Close) / series[0].Close) * 100,
			}

			// return the quotes
			resp <- q

		}(ticker)
	}

	wg.Wait()

	// check for errors
	if len(errCh) != 0 {
		return nil, <-errCh
	}

	// Create list of quotes from responses
	var quotes []*Quote
	for q := range resp {
		quotes = append(quotes, q)
	}

	return quotes, nil
}

// News returns recent news for a ticker NOTE: alphavantage doesn't have a news API, so use IEX instead
func (w *AlphaWrapper) News(ticker string) ([]string, error) {
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
