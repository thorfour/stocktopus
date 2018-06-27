package iex

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/thorfour/iex/pkg/endpoint"
	"github.com/thorfour/iex/pkg/types"
)

// Quote returns a stock quote for a given ticker
func Quote(ticker string) (*types.Quote, error) {
	url := endpoint.API().Stock().Ticker(ticker).Quote()
	jsonQuote, err := getJSON(url)
	if err != nil {
		return nil, err
	}

	// Parse into quote
	var quote types.Quote
	err = json.Unmarshal(jsonQuote, &quote)
	if err != nil {
		return nil, err
	}

	return &quote, nil
}

// Price returns the current price of a ticker
func Price(ticker string) (float64, error) {
	url := endpoint.API().Stock().Ticker(ticker).Price()
	jsonQuote, err := getJSON(url)
	if err != nil {
		return -1, err
	}

	price, err := strconv.ParseFloat(string(jsonQuote), 64)
	if err != nil {
		return -1, err
	}

	return price, nil
}

// BatchQuotes returns quotes for multiple tickers using a batch request
func BatchQuotes(tickers []string) (types.Batch, error) {

	url := endpoint.API().Stock().Market().Batch().Symbols().Tickers(tickers).And().Types(types.QuoteStr)
	jsonQuote, err := getJSON(url)
	if err != nil {
		return nil, err
	}

	// Parse into quote
	var batch types.Batch
	err = json.Unmarshal(jsonQuote, &batch)
	if err != nil {
		return nil, err
	}

	return batch, nil
}

// News returns the news for a given symbol
func News(ticker string) ([]types.News, error) {

	url := endpoint.API().Stock().Ticker(ticker).News().Last().Integer(5)
	jsonQuote, err := getJSON(url)
	if err != nil {
		return nil, err
	}

	// Parse into News
	var news []types.News
	err = json.Unmarshal(jsonQuote, &news)
	if err != nil {
		return nil, err
	}

	return news, nil
}

// getJSON returns the JSON response from a url
func getJSON(url endpoint.APIString) ([]byte, error) {

	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	// Read the quote into the slice
	defer resp.Body.Close()
	jsonQuote, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return jsonQuote, nil
}
