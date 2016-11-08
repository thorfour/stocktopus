//+build AWS

package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/colinmc/aws"
	"github.com/colinmc/stock"
)

const (
	addToList      = "WATCH"
	printList      = "LIST"
	removeFromList = "UNWATCH"
	clear          = "CLEAR"
)

type stockFunc func(string) (string, error)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
		return
	}

	// Expect args(1) to be a url encoded string
	decodedMap, err := url.ParseQuery(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: url.ParseQuery")
		return
	}

	text := decodedMap["text"]
	text = strings.Split(strings.ToUpper(text[0]), " ")

	switch text[0] {
	case addToList: // Add ticker to a watch list

		if len(text) < 2 { // Must be something to add to watch list
			fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
			return
		}

		// Chop off addToList arg
		text = text[1:]

		// User and token to be used as watch list lookup
		user := decodedMap["user_id"]
		token := decodedMap["token"]
		key := fmt.Sprintf("%v%v", token, user)

		err := aws.AddToList(key, text)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Error addtolist: %v", err))
			return
		}

		fmt.Print("Added")

	case removeFromList:

		if len(text) != 2 { // Only allow removal of 1 item
			fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
			return
		}

		// Chop off printList arg
		text = text[1:]

		// User and token to be used as watch list lookup
		user := decodedMap["user_id"]
		token := decodedMap["token"]
		key := fmt.Sprintf("%v%v", token, user)

		// Remove from watch list
		err := aws.RmFromList(key, text)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Error rmfromlist: %v", err))
			return
		}

		fmt.Print("Removed")

	case clear: // Remove entire watch list

		if len(text) > 1 {
			fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
			return
		}

		user := decodedMap["user_id"]
		token := decodedMap["token"]
		key := fmt.Sprintf("%v%v", token, user)

		err := aws.Clear(key)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Error clear: %v", err))
		}

		fmt.Print("Removed")

	case printList: // Print out all tickers in watch list

		if len(text) > 1 { // Requested more than just LIST
			fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
			return
		}

		// User and token to be used as watch list lookup
		user := decodedMap["user_id"]
		token := decodedMap["token"]
		key := fmt.Sprintf("%v%v", token, user)

		// Get and print watch list
		list, err := aws.GetList(key)
		if err != nil || len(list) == 0 {
			fmt.Fprintln(os.Stderr, "Error: No List")
			return
		}

		// Set the tickers to the list that was read. Fallthrough to normal printing
		text = strings.Split(list, " ")
		fallthrough

	default: // List of tickers to get information about right now

		var quoteFunc stockFunc
		var chartFunc stockFunc
		var quote string

		for _, ticker := range text {

			// Currently the longest stock ticker is 5 letters.
			// If a ticker is 6 characters assume a currency request
			if len(ticker) == 6 {
				quoteFunc = stock.GetCurrencyGoogle
			} else {
				quoteFunc = stock.GetQuoteGoogle
			}

			// Pull the quote
			q, err := quoteFunc(ticker)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error: ", err)
				return
			}

			quote = fmt.Sprintf("%v%v\n", quote, q)
		}

		// Pull a chart if single stock requested
		if len(text) == 1 {

			if len(text[0]) == 6 {
				chartFunc = stock.GetChartLinkCurrencyFinviz
			} else {
				chartFunc = stock.GetChartLinkFinviz
			}

			// Pull a stock chart
			chartUrl, err := chartFunc(text[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error: ", err)
				return
			}

			// Dump the chart link to stdio
			quote = fmt.Sprintf("%v%v", quote, chartUrl)
		}

		// Dump the quote to stdio
		fmt.Println(quote)
	}
}
