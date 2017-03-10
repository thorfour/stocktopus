package stock

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestQuoteGoogle(t *testing.T) {

	resp, err := GetQuoteGoogle("AMD")
	if err != nil {
		t.Fail()
	}
	fmt.Println(resp)
	resp, err = GetQuoteGoogle("TWLO")
	if err != nil {
		t.Fail()
	}
	fmt.Println(resp)
	resp, err = GetQuoteGoogle("WDC")
	if err != nil {
		t.Fail()
	}
	fmt.Println(resp)
}

func TestQuoteMOD(t *testing.T) {

	resp, err := GetQuoteMOD("AMD")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(resp)
}

func TestChartYahoo(t *testing.T) {

	resp, err := GetChartYahoo("MSFT")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Write resp to file
	ioutil.WriteFile("chart.png", resp, 0644)
}

func TestChartCompareGoogle(t *testing.T) {

	resp, err := GetChartLinkCompareGoogle("MSFT AAPL")
	if err != nil {
		t.Fail()
	}
	fmt.Println(resp)
}

func TestChartMOD(t *testing.T) {

	resp, err := GetChartJsonMOD("AAPL")
	if err != nil {
		t.Fail()
	}
	fmt.Println(resp)
}

func TestCurrencyGoogle(t *testing.T) {

	resp, err := GetCurrencyGoogle("GBPUSD")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(resp)
}

func TestGetInfo(t *testing.T) {

	tickers := []string{"GOOG", "ION", "COW"}
	for _, ticker := range tickers {
		resp, err := GetInfo(ticker)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}
		fmt.Println(resp)
	}
}

func TestGetBadInfo(t *testing.T) {

	tickers := []string{"TLND"}
	for _, ticker := range tickers {
		resp, err := GetInfo(ticker)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}
		fmt.Println(resp)
	}
}

func TestGetCurrencyYahoo(t *testing.T) {

	ticker := "BTCUSD"
	quote, err := GetCurrencyYahoo(ticker)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	fmt.Println(quote)
}

func TestGetPriceGoogleMulti(t *testing.T) {

	symbol := "GOOG AAPL MSFT TSLA"
	price, err := GetPriceGoogleMulti(symbol)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	fmt.Println(price)
}
