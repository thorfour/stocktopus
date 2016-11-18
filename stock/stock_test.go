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

	resp, err := GetInfo("GOOG")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(resp)
}
