package stock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Pulls a stock quote from google finance
// Assumes the format is passed back in json
func GetQuote(symbol string) string {

	symbol = strings.ToUpper(symbol)

	// Check for nasdaq first
	url := fmt.Sprintf("http://finance.google.com/finance/info?client=ig&q=%v", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("Failed to get quote: %v", err)
	}

	// Read the quote into the slice
	defer resp.Body.Close()
	jsonQuote, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Unable to read body: %v", err)
	}

	// Google quotes start with '//' as the response
	// as well as surrounding the json with '[]'
	jsonQuote = jsonQuote[6 : len(jsonQuote)-2]

	var q interface{}
	err = json.Unmarshal(jsonQuote, &q)
	if err != nil {
		return fmt.Sprintf("Unable to parse quote: %v, %v", err, string(jsonQuote))
	}

	// Type assertion
	quote, ok := q.(map[string]interface{})
	if !ok {
		return fmt.Sprintf("Quote was in unexpected format")
	}

	// Pull the current price and the change
	l_cur := quote["l_cur"]
	c := quote["c"]

	return fmt.Sprintf("%v Current Price: %v Todays Change: %v", symbol, l_cur, c)
}
