//+build AWS

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/colinmc/stock"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
		return
	}

	// Expected:  single arg with multiple tickers
	tickers := strings.Split(os.Args[1], " ")

	if len(tickers) == 1 {

		ticker := os.Args[1]

		// Pull the stock quote
		quote, err := stock.GetQuoteGoogle(ticker)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: ", err)
			return
		}
		fmt.Println(quote)

		// Pull a stock chart if only 1 ticker was sent
		if len(tickers) == 1 {
			chartUrl, err := stock.GetChartLinkFinviz(ticker)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error: ", err)
				return
			}
			fmt.Println(chartUrl)
		}

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
