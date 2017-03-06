package main

import "testing"

func TestGetQuotes(t *testing.T) {
	getQuotes([]string{"tsla", "amd", "wdc", "intc", "gpro", "f", "goog"}, nil)
}

func TestGetQuotesSingle(t *testing.T) {
	getQuotes([]string{"tsla"}, nil)
}

func TestGetQuotesSingleWithCurrency(t *testing.T) {
	getQuotes([]string{"tsla", "btcusd", "amd"}, nil)
}
