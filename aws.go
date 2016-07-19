//+build AWS

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/colinmc/stock"
)

func main() {

	// Expect
	// text - command text
	if len(os.Args) != 2 {
		return
	}

	ticker := os.Args[1]

	// Pull the stock quote
	quote, err := stock.GetQuote(ticker)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(quote)
}
