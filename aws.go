//+build !RTM

package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/bndr/gotabulate"
	"github.com/stocktopus/aws"
	"github.com/stocktopus/stock"
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

	// Play money commands
	buy       = "BUY"
	sell      = "SELL"
	short     = "SHORT"
	deposit   = "DEPOSIT"
	portfolio = "PORTFOLIO"
	reset     = "RESET"
)

var cmds map[string]cmdInfo

// Mapping of command string to function
func init() {
	cmds = map[string]cmdInfo{
		addToList:      {add, "*watch [tickers...]* add tickers to personal watch list"},
		printList:      {print, "*list*               print out personal watch list"},
		removeFromList: {remove, "*unwatch [ticker]*   remove single ticker from watch list"},
		clear:          {clearList, "*clear*              remove entire watch list"},
		info:           {getInfo, "*info [ticker]* print a company profile"},
		deposit:        {depositPlay, "*deposit [amount]* deposit amount of play money into account"},
		reset:          {resetPlay, "*reset* resets account"},
		portfolio:      {portfolioPlay, "*portfolio* Prints current portfolio of play money"},
		buy:            {buyPlay, "*buy [ticker shares]* Purchases number of shares in a security with play money"},
		sell:           {sellPlay, "*sell [ticker shares]* Sells number of shares of specified security"},
		help:           {printHelp, "*[tickers...]*       pull stock quotes for list of tickers"},
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

	// Accumulate all the quotes in the channel
	quotes := make([]string, len(text))
	wg := new(sync.WaitGroup)
	wg.Add(len(text))

	for i, ticker := range text {

		// Currently the longest stock ticker is 5 letters.
		// If a ticker is 6 characters assume a currency request
		if len(ticker) == 6 {
			quoteFunc = stock.GetCurrencyYahoo
		} else {
			quoteFunc = stock.GetQuoteGoogle
		}

		// Pull the quote
		go func(t string, index int) {
			q, err := quoteFunc(t)
			if err == nil {
				quotes[index] = q // Push the quote into the queue
			}
			wg.Done()
		}(text[i], i)
	}

	// Wait for all the quotes to complete
	wg.Wait()

	rows := make([][]string, len(quotes))
	for i, q := range quotes {
		info := strings.Fields(q)
		rows[i] = []string{info[0], info[3], info[6]}
	}

	t := gotabulate.Create(rows)
	t.SetHeaders([]string{"Ticker", "Current Price", "Todays Change"})
	t.SetAlign("left")
	t.SetHideLines([]string{"bottomLine", "betweenLine", "top"})
	quote = t.Render("simple")
	quote = fmt.Sprintf("```%v```", quote)

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

	if len(text) != 2 {
		fmt.Fprintln(os.Stderr, "Error: Invalid number of arguments")
		return
	}

	// Chop off arg
	text = text[1:]

	resp, err := stock.GetInfo(text[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		return
	}

	fmt.Println(resp)
}

//--------------------
// Play money
//--------------------

// account is a users play money account information
type account struct {
	Balance  float64
	Holdings map[string]Holding
}

type Holding struct {
	Strike float64 // Strike price of the purchase
	Shares uint64  // number of shares being held
}

func depositPlay(text []string, decodedMap url.Values) {

	if len(text) != 2 { // Must have an amount to add
		fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
		return
	}

	// Chop off deposit arg
	text = text[1:]

	// Parse amount to add to account
	amt, err := strconv.ParseUint(text[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Invalid amount: %v", err))
		return
	}

	// User and token to be used as lookup
	user := decodedMap["user_id"]
	token := decodedMap["token"]
	key := fmt.Sprintf("%v%v%v", "ACCT", token, user)

	// Load the account
	acct, err := loadAccount(key)
	if err != nil {
		// If no file exits then create a new account
		newAcct := new(account)
		newAcct.Holdings = make(map[string]Holding)
		newAcct.Balance = float64(amt)
		saveAccount(newAcct, key)
		fmt.Fprintln(os.Stderr, fmt.Sprintf("New Balance: %v", newAcct.Balance))
		return
	}

	// Add amount to balance
	acct.Balance += float64(amt)

	err = saveAccount(acct, key)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Unable to save account: %v", err))
	}

	// Respond with the new balance
	resp := fmt.Sprintf("New Balance: %v", acct.Balance)
	fmt.Fprintln(os.Stderr, resp)
}

func resetPlay(text []string, decodedMap url.Values) {

	if len(text) != 1 { // Only reset accepted
		fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
		return
	}

	// User and token to be used as lookup
	user := decodedMap["user_id"]
	token := decodedMap["token"]
	key := fmt.Sprintf("%v%v%v", "ACCT", token, user)

	newAcct := new(account)
	newAcct.Holdings = make(map[string]Holding)
	newAcct.Balance = float64(0)
	saveAccount(newAcct, key)
	fmt.Fprintln(os.Stderr, fmt.Sprintf("New Balance: %v", newAcct.Balance))
}

func portfolioPlay(text []string, decodedMap url.Values) {

	if len(text) != 1 { // Only portfolio accepted
		fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
		return
	}

	// User and token to be used as lookup
	user := decodedMap["user_id"]
	token := decodedMap["token"]
	key := fmt.Sprintf("%v%v%v", "ACCT", token, user)

	acct, err := loadAccount(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to load account")
		return
	}

	s := fmt.Sprintf("Balance: $%0.2f", acct.Balance)
	for k, v := range acct.Holdings {
		s = fmt.Sprintf("%v\n%v : %v @ $%v", s, k, v.Shares, v.Strike)
	}

	fmt.Print(s)
}

func buyPlay(text []string, decodedMap url.Values) {

	if len(text) != 3 { // BUY TICKER SHARES
		fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
		return
	}

	// Chop off buy arg
	text = text[1:]

	// Parse number of shares to purchase
	amt, err := strconv.ParseUint(text[1], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Invalid amount: %v", err))
		return
	}

	// lookup ticker price
	ticker := text[0]
	price, err := stock.GetPriceGoogle(ticker)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Unable to get price: %v", err))
		return
	}

	// User and token to be used as lookup
	user := decodedMap["user_id"]
	token := decodedMap["token"]
	key := fmt.Sprintf("%v%v%v", "ACCT", token, user)

	acct, err := loadAccount(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to load account")
		return
	}

	// check if enough money in account
	if acct.Balance < (price * float64(amt)) {
		fmt.Fprintln(os.Stderr, "Insufficient funds")
		return
	}

	// add to account
	acct.Balance -= (price * float64(amt))
	h, ok := acct.Holdings[ticker]
	if !ok {
		acct.Holdings[ticker] = Holding{price, amt}
	} else {
		newShares := h.Shares + amt
		acct.Holdings[ticker] = Holding{price, newShares}
	}

	// write account
	err = saveAccount(acct, key)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Unable to save account: %v", err))
		return
	}

	fmt.Fprintln(os.Stderr, "Done")
}

func sellPlay(text []string, decodedMap url.Values) {

	if len(text) != 3 { // SELL TICKER SHARES
		fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
		return
	}

	// Chop off buy arg
	text = text[1:]

	// Parse number of shares to purchase
	amt, err := strconv.ParseUint(text[1], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Invalid amount: %v", err))
		return
	}

	// lookup ticker price
	ticker := text[0]
	price, err := stock.GetPriceGoogle(ticker)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Unable to get price: %v", err))
		return
	}

	// User and token to be used as lookup
	user := decodedMap["user_id"]
	token := decodedMap["token"]
	key := fmt.Sprintf("%v%v%v", "ACCT", token, user)

	acct, err := loadAccount(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to load account")
		return
	}

	h, ok := acct.Holdings[ticker]
	if !ok || h.Shares < amt {
		fmt.Fprintln(os.Stderr, "Not enough shares")
		return
	}

	// remove from account and credit account for the sale
	newShares := h.Shares - amt
	if newShares == 0 {
		delete(acct.Holdings, ticker)
	} else {
		acct.Holdings[ticker] = Holding{h.Strike, newShares}
	}

	acct.Balance += float64(amt) * price

	// write account
	err = saveAccount(acct, key)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Unable to save account: %v", err))
		return
	}

	fmt.Fprintln(os.Stderr, "Done")
}

func saveAccount(acct *account, key string) error {

	// Encode account
	var eb bytes.Buffer
	e := gob.NewEncoder(&eb)
	err := e.Encode(acct)
	if err != nil {
		return fmt.Errorf("Unable to encode account")
	}

	// Save the account to file
	return aws.WriteFile(key, eb.Bytes())
}

func loadAccount(key string) (*account, error) {

	f, err := aws.LoadFile(key)
	if err != nil {
		return nil, fmt.Errorf("Unable to load account: %v", err)
	}

	// Unserialize acct from file
	acct := new(account)
	buf := bytes.NewBuffer(f)
	d := gob.NewDecoder(buf)
	err = d.Decode(acct)
	if err != nil {
		return nil, fmt.Errorf("Unable to decode account")
	}

	return acct, nil
}
