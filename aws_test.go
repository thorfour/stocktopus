package main

import (
	"fmt"
	"testing"
)

func TestGetQuotes(t *testing.T) {
	getQuotes([]string{"tsla", "amd", "wdc", "intc", "gpro", "f", "goog"}, nil)
}

func TestGetQuotesSingle(t *testing.T) {
	getQuotes([]string{"tsla"}, nil)
}

func TestGetQuotesSingleWithCurrency(t *testing.T) {
	getQuotes([]string{"tsla", "btcusd", "amd"}, nil)
}

func TestGetQuotesBad(t *testing.T) {
	fmt.Println("Middle -----------------------------------------")
	getQuotes([]string{"tsla", "osghoevcmi", "amd"}, nil)
	fmt.Println("Only-----------------------------------------")
	getQuotes([]string{"osghoevcmi"}, nil)
	fmt.Println("Start-----------------------------------------")
	getQuotes([]string{"osghoevcmi", "amd", "tsla"}, nil)
	fmt.Println("End-----------------------------------------")
	getQuotes([]string{"tsla", "amd", "aorghreqcm"}, nil)
	fmt.Println("Two-to-one-------------------------------------")
	getQuotes([]string{"amd", "aorghreqcm"}, nil)
}
