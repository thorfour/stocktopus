package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	redis "gopkg.in/redis.v5"

	"github.com/bndr/gotabulate"
	iex "github.com/thorfour/iex/pkg/api"
	iextypes "github.com/thorfour/iex/pkg/types"
	"github.com/thorfour/stocktopus/pkg/cfg"
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
	news           = "NEWS"

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
		news:           {getNews, "*ticker* Displays the latest news for a company"},
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
		getQuotes(decodedMap["text"][0], decodedMap)
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

	// If the first arg starts with '#' then it's the name of the list
	if text[0][0] == '#' {
		user = []string{strings.ToLower(text[0][1:]), decodedMap["team_id"][0]}
		text = text[1:] // Remove list name
	}

	key := fmt.Sprintf("%v%v", token, user)

	rClient := connectRedis()

	// Convert []string to []interface{} for the SAdd call
	members := []interface{}{}
	for _, member := range text {
		members = append(members, interface{}(member))
	}

	_, err := rClient.SAdd(key, members...).Result()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Error addtolist: %v", err))
		return
	}

	fmt.Fprintln(os.Stderr, "Added")
}

// Print out a watchlist
func print(text []string, decodedMap url.Values) {

	// User and token to be used as watch list lookup
	user := decodedMap["user_id"]
	token := decodedMap["token"]

	// Chop off printList arg
	text = text[1:]

	// If the first arg starts with '#' then it's the name of the list
	if len(text) == 1 && text[0][0] == '#' {
		user = []string{strings.ToLower(text[0][1:]), decodedMap["team_id"][0]}
		text = text[1:] // Remove list name
	} else if len(text) >= 1 {
		fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
		return
	}

	key := fmt.Sprintf("%v%v", token, user)

	rClient := connectRedis()

	// Get and print watch list
	list, err := rClient.SMembers(key).Result()
	if err != nil || len(list) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No List")
		return
	}

	getQuotes(strings.Join(list, " "), decodedMap)
}

// Remove a single ticker from a watch list
func remove(text []string, decodedMap url.Values) {

	// Chop off printList arg
	text = text[1:]

	// User and token to be used as watch list lookup
	user := decodedMap["user_id"]
	token := decodedMap["token"]

	// If the first arg starts with '#' then it's the name of the list
	if len(text) > 1 && text[0][0] == '#' {
		user = []string{strings.ToLower(text[0][1:]), decodedMap["team_id"][0]}
		text = text[1:] // Remove list name
	}

	key := fmt.Sprintf("%v%v", token, user)

	rClient := connectRedis()

	// Convert []string to []interface{} for the SRem call
	members := []interface{}{}
	for _, member := range text {
		members = append(members, interface{}(member))
	}

	// Remove from watch list
	_, err := rClient.SRem(key, members...).Result()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Error rmfromlist: %v", err))
		return
	}

	fmt.Fprintln(os.Stderr, "Removed")
}

// Delete a watch list. Deletes the whole file instead of clearing
func clearList(text []string, decodedMap url.Values) {

	user := decodedMap["user_id"]
	token := decodedMap["token"]

	// Chop off printList arg
	text = text[1:]

	// If the first arg starts with '#' then it's the name of the list
	if len(text) == 1 && text[0][0] == '#' {
		user = []string{strings.ToLower(text[0][1:]), decodedMap["team_id"][0]}
		text = text[1:] // Remove list name
	} else if len(text) >= 1 {
		fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
		return
	}

	key := fmt.Sprintf("%v%v", token, user)

	rClient := connectRedis()

	_, err := rClient.Del(key).Result()
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

// text is expected to be a list of tickers separated by spaces
func getMultiQuote(text string) (iextypes.Batch, error) {
	tickers := strings.Split(text, " ")
	batch, err := iex.BatchQuotes(tickers)
	if err != nil {
		return nil, err
	}

	return batch, nil
}

// Default functionality of grabbing stock quote(s)
func getQuotes(text string, decodedMap url.Values) {
	var chartFunc stockFunc
	var quote string

	// Pull the quote
	info, err := getMultiQuote(text)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		return
	}

	// Nothing was returned
	if len(info) == 0 {
		fmt.Fprintln(os.Stderr, "There's nothing here")
		return
	}

	rows := make([][]interface{}, 0, len(info))
	cumsum := float64(0)
	for ticker := range info {
		quote, err := info.Quote(ticker)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to get quote for ", ticker)
			return
		}
		rows = append(rows, []interface{}{ticker, quote.IexRealtimePrice, quote.Change, fmt.Sprintf("%0.3f", (100 * quote.ChangePercent))})
		cumsum += (100 * quote.ChangePercent)
	}
	rows = append(rows, []interface{}{"Avg.", "---", "---", fmt.Sprintf("%0.3f%%", cumsum/float64(len(rows)))})

	t := gotabulate.Create(rows)
	t.SetHeaders([]string{"Company", "Current Price", "Todays Change", "Percent Change"})
	t.SetAlign("right")
	t.SetHideLines([]string{"bottomLine", "betweenLine", "top"})
	quote = t.Render("simple")
	quote = fmt.Sprintf("```%v```", quote)

	// Pull a chart if single stock requested
	if len(rows) == 2 {

		if len(text) == 6 {
			chartFunc = getChartLinkCurrencyFinviz
		} else {
			chartFunc = getChartLinkFinviz
		}

		// Pull a stock chart
		chartURL, err := chartFunc(text)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: ", err)
			return
		}

		// Dump the chart link to stdio
		quote = fmt.Sprintf("%v\n%v", quote, chartURL)
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

	resp, err := getStockInfo(text[0])
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
	Holdings map[string]holding
}

type holding struct {
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

	client := connectRedis()

	// Load the account
	acct, err := loadAccount(client, key)
	if err != nil {
		// If no file exits then create a new account
		newAcct := new(account)
		newAcct.Holdings = make(map[string]holding)
		newAcct.Balance = float64(amt)
		saveAccount(client, newAcct, key)
		fmt.Fprintln(os.Stderr, fmt.Sprintf("New Balance: %v", newAcct.Balance))
		return
	}

	// Add amount to balance
	acct.Balance += float64(amt)

	err = saveAccount(client, acct, key)
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

	client := connectRedis()

	newAcct := new(account)
	newAcct.Holdings = make(map[string]holding)
	newAcct.Balance = float64(0)
	saveAccount(client, newAcct, key)
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

	client := connectRedis()

	acct, err := loadAccount(client, key)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to load account")
		return
	}

	total := float64(0)
	totalChange := float64(0)
	if len(acct.Holdings) > 0 {

		var list []string // List of all tickers
		for ticker := range acct.Holdings {
			list = append(list, ticker)
		}

		// Pull the quote
		info, err := getMultiQuote(strings.Join(list, " "))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to get quotes")
			return
		}

		rows := make([][]interface{}, 0, len(acct.Holdings))
		for ticker := range info {
			quote, err := info.Quote(ticker)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Unable to get quote for ", ticker)
				return
			}
			h := acct.Holdings[ticker]
			total += float64(h.Shares) * quote.IexRealtimePrice
			delta := float64(h.Shares) * (quote.IexRealtimePrice - h.Strike)
			totalChange += delta
			deltaStr := fmt.Sprintf("%0.2f", delta)
			rows = append(rows, []interface{}{ticker, h.Shares, h.Strike, quote.IexRealtimePrice, deltaStr})
		}

		rows = append(rows, []interface{}{"Total", "---", "---", "---", fmt.Sprintf("%0.2f", totalChange)})

		t := gotabulate.Create(rows)
		t.SetHeaders([]string{"Ticker", "Shares", "Strike", "Current", "Gain/Loss $"})
		t.SetAlign("left")
		t.SetHideLines([]string{"bottomLine", "betweenLine", "top"})
		table := t.Render("simple")
		summary := fmt.Sprintf("Portfolio Value: $%0.2f\nBalance: $%0.2f\nTotal: $%0.2f", total, acct.Balance, total+acct.Balance)
		resp := fmt.Sprintf("```%v\n%v```", table, summary)
		fmt.Print(resp)
	} else {
		s := fmt.Sprintf("Balance: $%0.2f", acct.Balance)
		fmt.Print(s)
	}
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
	price, err := iex.Price(ticker)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Unable to get price: %v", err))
		return
	}

	// User and token to be used as lookup
	user := decodedMap["user_id"]
	token := decodedMap["token"]
	key := fmt.Sprintf("%v%v%v", "ACCT", token, user)

	client := connectRedis()

	acct, err := loadAccount(client, key)
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
		acct.Holdings[ticker] = holding{price, amt}
	} else {
		newShares := h.Shares + amt
		acct.Holdings[ticker] = holding{price, newShares}
	}

	// write account
	err = saveAccount(client, acct, key)
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
	price, err := iex.Price(ticker)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Unable to get price: %v", err))
		return
	}

	// User and token to be used as lookup
	user := decodedMap["user_id"]
	token := decodedMap["token"]
	key := fmt.Sprintf("%v%v%v", "ACCT", token, user)

	client := connectRedis()

	acct, err := loadAccount(client, key)
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
		acct.Holdings[ticker] = holding{h.Strike, newShares}
	}

	acct.Balance += float64(amt) * price

	// write account
	err = saveAccount(client, acct, key)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Unable to save account: %v", err))
		return
	}

	fmt.Fprintln(os.Stderr, "Done")
}

func saveAccount(client *redis.Client, acct *account, key string) error {

	// Encode account
	serialized, err := json.Marshal(acct)
	if err != nil {
		return fmt.Errorf("Unable to encode account")
	}

	// Save the account to redis
	_, err = client.Set(key, string(serialized), 0).Result()

	return err
}

func loadAccount(client *redis.Client, key string) (*account, error) {

	serialized, err := client.Get(key).Result()
	if err != nil {
		return nil, fmt.Errorf("Unable to load account: %v", err)
	}

	// Unserialize acct from file
	acct := new(account)
	err = json.Unmarshal([]byte(serialized), acct)
	if err != nil {
		return nil, fmt.Errorf("Unable to decode account")
	}

	return acct, nil
}

func connectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPw,
		DB:       0,
	})
}

// getNews will print news from requested company
func getNews(text []string, decodedMap url.Values) {
	if len(text) != 2 {
		fmt.Fprintln(os.Stderr, "Error: Invalid number of arguments")
	}
	// Chop off news arg
	text = text[1:]

	latestNews, err := iex.News(text[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: No news is good news right?")
	}

	// Try and pretty print them
	var printNews string
	for _, n := range latestNews {
		printNews = fmt.Sprintf("%s%s\n\n", printNews, n.Summary)
	}

	fmt.Println(printNews)
}

// getStockInfo returns a company information paragrah from reuters
func getStockInfo(symbol string) (string, error) {

	symbol = strings.ToUpper(symbol)
	url := fmt.Sprintf("http://reuters.com/finance/stocks/companyProfile?symbol=%v", symbol)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	tokenizer := html.NewTokenizer(resp.Body)
	nextParagraph := false
	moduleBody := false
	for {
		token := tokenizer.Next()
		if token == html.ErrorToken {
			break
		}

		if token != html.StartTagToken {
			if nextParagraph {
				text := string(tokenizer.Text())
				if len(text) > 3 {
					return text, nil
				}
			}
			continue
		}

		switch {
		case moduleBody:
			tag, hasAttr := tokenizer.TagName()
			if string(tag) == "div" && hasAttr {
				key, val, _ := tokenizer.TagAttr()
				if string(key) == "class" && string(val) == "moduleBody" {
					nextParagraph = true
					moduleBody = false
				}
			}
		case nextParagraph:
			tag, _ := tokenizer.TagName()
			switch string(tag) {
			case "p":
				tokenizer.Next()
				t := string(tokenizer.Text())
				if len(t) == 0 {
					t = "There's nothing here"
				}
				return t, nil
			}
		default:
			// Find <div id="companyNews">
			// after that the following tag to look for is <div class="moduleBody">
			tag, hasAttr := tokenizer.TagName()
			if string(tag) == "div" && hasAttr {
				key, val, _ := tokenizer.TagAttr()
				if string(key) == "id" && string(val) == "companyNews" {
					moduleBody = true
				}
			}
		}
	}

	return "", fmt.Errorf("Unable to find quote")
}

// getChartLinkCurrencyFinviz returns a currenct chart link from finviz
func getChartLinkCurrencyFinviz(symbol string) (string, error) {

	symbol = strings.ToUpper(symbol)
	url := fmt.Sprintf("http://finviz.com/fx_image.ashx?%v_d1_l.png", symbol)

	return url, nil
}

// getChartLinkFinviz returns a chart link from finviz
func getChartLinkFinviz(symbol string) (string, error) {

	symbol = strings.ToUpper(symbol)
	url := fmt.Sprintf("http://finviz.com/chart.ashx?t=%v&ty=c&ta=1&p=d&s=l", symbol)

	return url, nil
}
