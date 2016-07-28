//+build AWS

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/colinmc/stock"
)

type stockFunc func(string) (string, error)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
		return
	}

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
