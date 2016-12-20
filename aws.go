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

type cmdFunc func([]string, url.Values)

type cmdInfo struct {
	funcPtr cmdFunc // Function pointer to the function to execute
	helpStr string  // help string
}

// Supported commands
const (
	addToList      = "WATCH"
	printList      = "LIST"
	removeFromList = "UNWATCH"
	clear          = "CLEAR"
	help           = "HELP"
	info           = "INFO"
)

var cmds map[string]cmdInfo

// Mapping of command string to function
func init() {
	cmds = map[string]cmdInfo{
		addToList:      cmdInfo{add, "*watch [tickers...]* add tickers to personal watch list"},
		printList:      cmdInfo{print, "*list*               print out personal watch list"},
		removeFromList: cmdInfo{remove, "*unwatch [ticker]*   remove single ticker from watch list"},
		clear:          cmdInfo{clearList, "*clear*              remove entire watch list"},
		info:           cmdInfo{getInfo, "*info [ticker]* print a company profile"},
		help:           cmdInfo{printHelp, "*[tickers...]*       pull stock quotes for list of tickers"},
	}
}

type stockFunc func(string) (string, error)

// Successful command print to stdout, errors and ephermeral messages print to stderr
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

	cmd, ok := cmds[text[0]]
	if !ok { // If there is no cmd mapped, assume it's a ticker and get quotes
		getQuotes(text, decodedMap)
	} else {
		cmd.funcPtr(text, decodedMap)
	}
}

// Add ticker(s) to a watch list
func add(text []string, decodedMap url.Values) {

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

	fmt.Fprintln(os.Stderr, "Added")
}

// Print out a watchlist
func print(text []string, decodedMap url.Values) {

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

	getQuotes(text, decodedMap)
}

// Remove a single ticker from a watch list
func remove(text []string, decodedMap url.Values) {

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

	fmt.Fprintln(os.Stderr, "Removed")
}

// Delete a watch list. Deletes the whole file instead of clearing
func clearList(text []string, decodedMap url.Values) {

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

	fmt.Fprintln(os.Stderr, "Removed")
}

// Prints out help information about supported commands
func printHelp(text []string, decodedMap url.Values) {

	var out string
	for _, val := range cmds {
		out = fmt.Sprintf("%v\n%v", out, val.helpStr)
	}

	fmt.Fprintln(os.Stderr, out)
}

// Default functionality of grabbing stock quote(s)
func getQuotes(text []string, decodedMap url.Values) {
	var quoteFunc stockFunc
	var chartFunc stockFunc
	var quote string

	for _, ticker := range text {

		// Currently the longest stock ticker is 5 letters.
		// If a ticker is 6 characters assume a currency request
		if len(ticker) == 6 {
			quoteFunc = stock.GetCurrencyYahoo
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

// Print out a company profile
func getInfo(text []string, decodedMap url.Values) {

	// Chop off arg
	text = text[1:]

	if len(text) > 1 {
		fmt.Fprintln(os.Stderr, "Error: Too many arguments")
		return
	}

	resp, err := stock.GetInfo(text[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		return
	}

	fmt.Println(resp)
}
