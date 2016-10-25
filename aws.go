//+build AWS

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/colinmc/stock"
)

const (
	addToList = "WATCH"
	printList = "LIST"
)

type stockFunc func(string) (string, error)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
		return
	}

	switch os.Args[1] {
	case addToList: // Add ticker to a watch list
		// TODO add to watch list
	case printList: // Print out all tickers in watch list
		// TODO print watch list
	default: // List of tickers to get information about right now

		// Expected:  single arg with multiple tickers
		tickers := strings.Split(os.Args[1], " ")

		if len(tickers) == 1 {

			ticker := os.Args[1]

			var quoteFunc stockFunc
			var chartFunc stockFunc

			// Currently the longest stock ticker is 5 letters.
			// If a ticker is 6 characters assume a currency request
			if len(ticker) == 6 {
				quoteFunc = stock.GetCurrencyGoogle
				chartFunc = stock.GetChartLinkCurrencyFinviz
			} else {
				quoteFunc = stock.GetQuoteGoogle
				chartFunc = stock.GetChartLinkFinviz
			}

			// Pull the quote
			quote, err := quoteFunc(ticker)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error: ", err)
				return
			}

			// Dump the quote to stdio
			fmt.Println(quote)

			// Pull a stock chart
			chartUrl, err := chartFunc(ticker)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error: ", err)
				return
			}

			// Dump the chart link to stdio
			fmt.Println(chartUrl)

		} else {

			// Pull a comparison chart
			chartUrl, err := stock.GetChartLinkCompareGoogle(os.Args[1])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error: ", err)
				return
			}
			fmt.Println(chartUrl)
		}
	}
}
