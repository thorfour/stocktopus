package stock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// Pulls a stock quote from google finance
// Assumes the format is passed back in json
func GetQuoteGoogle(symbol string) (string, error) {

	symbol = strings.ToUpper(symbol)

	url := fmt.Sprintf("http://finance.google.com/finance/info?client=ig&q=%v", symbol)

	return parseGoogleFinanceResp(url)
}

func GetCurrencyGoogle(symbol string) (string, error) {

	symbol = strings.ToUpper(symbol)

	url := fmt.Sprintf("http://finance.google.com/finance/info?q=CURRENCY:%v", symbol)

	return parseGoogleFinanceResp(url)
}

func parseGoogleFinanceResp(url string) (string, error) {

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
	t := quote["t"]
	l_cur := quote["l_cur"]
	c := quote["c"]
	cp := quote["cp"]

	return fmt.Sprintf("*%v*\tCurrent Price: %v\tTodays Change: %v(%v%%)", t, l_cur, c, cp), nil
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

func GetChartJsonMOD(symbol string) (string, error) {

	type data struct {
		Min     float64   `json:"min"`
		Max     float64   `json:"max"`
		MaxDate string    `json:"maxDate"`
		MinDate string    `json:"minDate"`
		Value   []float64 `json:"values"`
	}

	type dataSeries struct {
		Close data `json:"close"`
	}

	type element struct {
		Currency   string     `json:"Currency"`
		TimeStamp  string     `json:"TimeStamp"`
		Symbol     string     `json:"Symbol"`
		Type       string     `json:"Type"`
		DataSeries dataSeries `json:"DataSeries"`
	}

	type MODChart struct {
		Positions []float64 `json:"Positions"`
		Dates     []string  `json:"Dates"`
		Elements  []element `json:"Elements"`
	}

	fileName := "options.js"
	symbol = strings.ToUpper(symbol)

	params := fmt.Sprintf("{\"Normalized\":\"false\",\"DataPeriod\":\"Day\",\"NumberOfDays\":\"365\",\"Elements\":[{\"Symbol\":\"%v\",\"Type\":\"price\",\"Params\":[\"c\"]}]}", symbol)
	url := fmt.Sprintf("http://dev.markitondemand.com/Api/v2/InteractiveChart/json?parameters=%v", params)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	// Read the quote into the slice
	defer resp.Body.Close()
	jsonChart, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Unmarshal into MODChart type
	var chart MODChart
	err = json.Unmarshal(jsonChart, &chart)
	if err != nil {
		return "", err
	}

	values := strings.Replace(fmt.Sprintf("%v", chart.Elements[0].DataSeries.Close.Value), " ", ",", -1)
	chartJson := fmt.Sprintf("{ chart: { type: 'line' }, title: { text: '%v' }, series: [{ name: '%v', data: %v}]}", symbol, symbol, values)

	err = ioutil.WriteFile(fileName, []byte(chartJson), 0644)
	if err != nil {
		return "", err
	}

	return fileName, nil
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
	url := fmt.Sprintf("http://finviz.com/chart.ashx?t=%v&ty=c&ta=1&p=d&s=l", symbol)

	return url, nil
}

func GetChartLinkCompareGoogle(symbols string) (string, error) {

	symbols = strings.ToUpper(symbols)

	// Replace spaces with commas for the chart url
	symbols = strings.Replace(symbols, " ", ",", -1)
	url := fmt.Sprintf("https://www.google.com/finance/chart?cht=c&q=%v&tlf=12h", symbols)

	return url, nil
}

func GetChartLinkCurrencyFinviz(symbol string) (string, error) {

	symbol = strings.ToUpper(symbol)
	url := fmt.Sprintf("http://finviz.com/fx_image.ashx?%v_d1_l.png", symbol)

	return url, nil
}

func GetInfo(symbol string) (string, error) {

	symbol = strings.ToUpper(symbol)
	url := fmt.Sprintf("http://reuters.com/finance/stocks/companyProfile?symbol=%v", symbol)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	tokenizer := html.NewTokenizer(resp.Body)
	nextParagraph := false
	for {
		token := tokenizer.Next()
		if token == html.ErrorToken {
			break
		}

		if nextParagraph {
			if token == html.StartTagToken {
				tag, _ := tokenizer.TagName()
				if string(tag) == "p" {
					tokenizer.Next()
					return string(tokenizer.Text()), nil
				}
			}

		} else {
			// Find <div id="companyNews">
			// after that the following tag to look for is <p>
			if token == html.StartTagToken {
				tag, hasAttr := tokenizer.TagName()
				if string(tag) == "div" && hasAttr {
					key, val, _ := tokenizer.TagAttr()
					if string(key) == "id" && string(val) == "companyNews" {
						nextParagraph = true
					}
				}
			}
		}
	}

	return "", fmt.Errorf("Unable to find quote")
}
