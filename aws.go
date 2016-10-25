//+build AWS

package main

import (
	"fmt"
	"os"

	"github.com/colinmc/stock"
)

const (
	addToList      = "WATCH"
	printList      = "LIST"
	removeFromList = "UNWATCH"
)

type stockFunc func(string) (string, error)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
		return
	}

	// Expected:  single json string
	// TODO THOR

	switch tickers[0] {
	case addToList: // Add ticker to a watch list

		if len(tickers) < 2 { // Must be something to add to watch list
			fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
			return
		}

		// Chop off addToList arg
		tickers = tickers[1:]

		// TODO add to watch list need to also obtain username etc...
		//aws.Watch(tickers)

	case printList: // Print out all tickers in watch list

		if len(tickers) > 1 { // Requested more than just LIST
			fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
			return
		}

		// Chop off printList arg
		tickers = tickers[1:]

		// TODO print watch list

	case removeFromList:

		// TODO remove from watch list

	default: // List of tickers to get information about right now

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
