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
		fmt.Println("Invalid arguments")
		return
	}

	// Expected:  single arg with multiple tickers
	tickers := strings.Split(os.Args[1], " ")

	for _, ticker := range tickers {
		// Pull the stock quote
		quote, err := stock.GetQuoteGoogle(ticker)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println(quote)

		// Pull a stock chart if only 1 ticker was sent
		if len(tickers) == 1 {
			chartUrl, err := stock.GetChartLinkFinviz(ticker)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Println(chartUrl)
		}
	}
}
