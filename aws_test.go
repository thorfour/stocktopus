package main

import (
	"fmt"
	"testing"
)

func TestGetQuotes(t *testing.T) {
	getQuotes("tsla amd wdc intc gpro f goog", nil)
}

func TestGetQuotesSingle(t *testing.T) {
	getQuotes("tsla", nil)
}

func TestGetQuotesSingleWithCurrency(t *testing.T) {
	getQuotes("tsla btcusd amd", nil)
}

func TestGetQuotesBad(t *testing.T) {
	fmt.Println("Middle -----------------------------------------")
	getQuotes("tsla osghoevcmi amd", nil)
	fmt.Println("Only-----------------------------------------")
	getQuotes("osghoevcmi", nil)
	fmt.Println("Start-----------------------------------------")
	getQuotes("osghoevcmi amd tsla", nil)
	fmt.Println("End-----------------------------------------")
	getQuotes("tsla amd aorghreqcm", nil)
	fmt.Println("Two-to-one-------------------------------------")
	getQuotes("amd aorghreqcm", nil)
}
