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
		quote, err := stock.GetQuote(ticker)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println(quote)
	}
}
