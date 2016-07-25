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
func GetQuoteGoogle(symbol string) (string, error) {

	symbol = strings.ToUpper(symbol)

	url := fmt.Sprintf("http://finance.google.com/finance/info?client=ig&q=%v", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	// Read the quote into the slice
	defer resp.Body.Close()
	jsonQuote, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Google quotes start with '//' as the response
	// as well as surrounding the json with '[]'
	jsonQuote = jsonQuote[6 : len(jsonQuote)-2]

	var q interface{}
	err = json.Unmarshal(jsonQuote, &q)
	if err != nil {
		return "", err
	}

	// Type assertion
	quote, ok := q.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf(fmt.Sprintf("Quote was in unexpected format"))
	}

	// Pull the current price and the change
	l_cur := quote["l_cur"]
	c := quote["c"]

	return fmt.Sprintf("*%v*\tCurrent Price: %v\tTodays Change: %v", symbol, l_cur, c), nil
}

// Pulls a stock quote from markit on demand
// markitondemand.com
func GetQuoteMOD(symbol string) (string, error) {

	symbol = strings.ToUpper(symbol)

	url := fmt.Sprintf("http://dev.markitondemand.com/Api/v2/Quote/json?symbol=%v", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	// Read the quote into the slice
	defer resp.Body.Close()
	jsonQuote, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var q interface{}
	err = json.Unmarshal(jsonQuote, &q)
	if err != nil {
		return "", err
	}

	// Type assertion
	quote, ok := q.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf(fmt.Sprintf("Quote was in unexpected format"))
	}

	// Pull the current price and the change
	l_cur := quote["LastPrice"]
	c := quote["Change"]

	return fmt.Sprintf("*%v*\tCurrent Price: %v\tTodays Change: %v", symbol, l_cur, c), nil
}

// Pulls a png stock image from yahoo finance
func GetChartYahoo(symbol string) ([]byte, error) {

	symbol = strings.ToUpper(symbol)

	url := fmt.Sprintf("http://chart.finance.yahoo.com/z?s=%v&t=6m&q=l&l=on&z=s&p=m50,m200", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// Read the quote into the slice
	defer resp.Body.Close()
	image, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func GetChartLinkYahoo(symbol string) (string, error) {

	symbol = strings.ToUpper(symbol)
	url := fmt.Sprintf("http://chart.finance.yahoo.com/z?s=%v&t=6m&q=l&l=on&z=s&p=m50,m200", symbol)

	return url, nil
}

func GetChartLinkFinviz(symbol string) (string, error) {

	symbol = strings.ToUpper(symbol)
	url := fmt.Sprintf("http://finviz.com/chart.ashx?t=%v&ty=c&ta=1&p=d&s=m", symbol)

	return url, nil
}

func GetChartLinkCompareGoogle(symbols string) (string, error) {

	// Replace spaces with commas for the chart url
	symbols = strings.Replace(symbols, " ", ",", -1)
	url := fmt.Sprintf("https://www.google.com/finance/chart?cht=c&q=%v&tlf=12h", symbols)

	return url, nil
}
