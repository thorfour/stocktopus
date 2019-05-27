package stocktopus

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	redis "gopkg.in/redis.v5"

	"github.com/bndr/gotabulate"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/thorfour/stocktopus/pkg/cfg"
	"github.com/thorfour/stocktopus/pkg/stock"
)

type cmdFunc func([]string, url.Values) (string, error)

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
	stats          = "STATS"

	// Play money commands
	buy       = "BUY"
	sell      = "SELL"
	short     = "SHORT"
	deposit   = "DEPOSIT"
	portfolio = "PORTFOLIO"
	reset     = "RESET"
)

var (
	cmds    map[string]cmdInfo
	cmdHist = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "command_timings",
		Help: "A histogram of cmd request execution times",
	}, []string{"command"})
)

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
		stats:          {getStats, "stats *ticker* [field options...]"},
		help:           {printHelp, "*[tickers...]*       pull stock quotes for list of tickers"},
	}
}

type stockFunc func(string) (string, error)

// measureTime is a helper function to measure the execution time of a function
func measureTime(start time.Time, label string) {
	cmdHist.WithLabelValues(label).Observe(time.Since(start).Seconds())
}

// Process url string to provide stocktpus functionality
func Process(args url.Values) (string, error) {
	text, ok := args["text"]
	if !ok {
		return "", errors.New("Bad request")
	}

	text = strings.Split(strings.ToUpper(text[0]), " ")
	cmd, ok := cmds[text[0]]
	if !ok {
		return getQuotes(args["text"][0], args)
	}
	return cmd.funcPtr(text, args)
}

// Add ticker(s) to a watch list
func add(text []string, decodedMap url.Values) (string, error) {
	defer measureTime(time.Now(), "add")

	if len(text) < 2 { // Must be something to add to watch list
		return "", errors.New("Error: Invalid number of arguments")
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
		return "", fmt.Errorf("Error addtolist: %v", err)
	}

	// Not an error but message of Added should be supressed
	return "", errors.New("Added")
}

// Print out a watchlist
func print(text []string, decodedMap url.Values) (string, error) {
	defer measureTime(time.Now(), "print")

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
		return "", errors.New("Error: Invalid number arguments")
	}

	key := fmt.Sprintf("%v%v", token, user)

	rClient := connectRedis()

	// Get and print watch list
	list, err := rClient.SMembers(key).Result()
	if err != nil || len(list) == 0 {
		return "", errors.New("Error: No List")
	}

	return getQuotes(strings.Join(list, " "), decodedMap)
}

// Remove a single ticker from a watch list
func remove(text []string, decodedMap url.Values) (string, error) {
	defer measureTime(time.Now(), "remove")

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
		return "", fmt.Errorf("Error rmfromlist: %v", err)
	}

	return "", errors.New("Removed")
}

// Delete a watch list. Deletes the whole file instead of clearing
func clearList(text []string, decodedMap url.Values) (string, error) {
	defer measureTime(time.Now(), "clear")

	user := decodedMap["user_id"]
	token := decodedMap["token"]

	// Chop off printList arg
	text = text[1:]

	// If the first arg starts with '#' then it's the name of the list
	if len(text) == 1 && text[0][0] == '#' {
		user = []string{strings.ToLower(text[0][1:]), decodedMap["team_id"][0]}
		text = text[1:] // Remove list name
	} else if len(text) >= 1 {
		return "", errors.New("Error: Invalid number arguments")
	}

	key := fmt.Sprintf("%v%v", token, user)

	rClient := connectRedis()

	_, err := rClient.Del(key).Result()
	if err != nil {
		return "", fmt.Errorf("Error clear: %v", err)
	}

	return "", errors.New("Removed")
}

// Prints out help information about supported commands
func printHelp(text []string, decodedMap url.Values) (string, error) {
	defer measureTime(time.Now(), "help")

	var out string
	for _, val := range cmds {
		out = fmt.Sprintf("%v\n%v", out, val.helpStr)
	}

	return "", errors.New(out)
}

// text is expected to be a list of tickers separated by spaces
func getMultiQuote(text string) ([]*stock.Quote, error) {
	tickers := strings.Split(text, " ")
	batch, err := stockInterface.BatchQuotes(tickers)
	if err != nil {
		return nil, err
	}

	return batch, nil
}

// Default functionality of grabbing stock quote(s)
func getQuotes(text string, decodedMap url.Values) (string, error) {
	defer measureTime(time.Now(), "quotes")

	var chartFunc stockFunc
	var quote string

	// Pull the quote
	info, err := getMultiQuote(text)
	if err != nil {
		return "", fmt.Errorf("Error: %v", err)
	}

	// Nothing was returned
	if len(info) == 0 {
		return "", errors.New("There's nothing here")
	}

	// sort info by changePercent
	sort.Sort(sortableList(info))

	rows := make([][]interface{}, 0, len(info))
	cumsum := float64(0)
	for _, quote := range info {
		rows = append(rows, []interface{}{quote.Ticker, quote.LatestPrice, fmt.Sprintf("%0.2f", quote.Change), fmt.Sprintf("%0.3f", (100 * quote.ChangePercent))})
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
			return "", fmt.Errorf("Error: %v", err)
		}

		quote = fmt.Sprintf("%v\n%v", quote, chartURL)
	}

	return quote, nil
}

// Print out a company profile
func getInfo(text []string, decodedMap url.Values) (string, error) {
	defer measureTime(time.Now(), "info")

	if len(text) != 2 {
		return "", errors.New("Error: Invalid number of arguments")
	}

	// Chop off arg
	text = text[1:]

	info, err := stockInterface.Company(text[0])
	if err != nil {
		return "", errors.New("Error: Unable to find info")
	}

	return strings.Join([]string{info.CompanyName, info.Industry, info.Website, info.CEO, info.Description}, "\n"), nil
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

func depositPlay(text []string, decodedMap url.Values) (string, error) {
	defer measureTime(time.Now(), "deposit")

	if len(text) != 2 { // Must have an amount to add
		return "", errors.New("Error: Invalid number arguments")
	}

	// Chop off deposit arg
	text = text[1:]

	// Parse amount to add to account
	amt, err := strconv.ParseUint(text[0], 10, 64)
	if err != nil {
		return "", fmt.Errorf("Invalid amount: %v", err)
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
		return "", fmt.Errorf("New Balance: %v", newAcct.Balance)
	}

	// Add amount to balance
	acct.Balance += float64(amt)

	err = saveAccount(client, acct, key)
	if err != nil {
		return "", fmt.Errorf("Unable to save account: %v", err)
	}

	// Respond with the new balance
	resp := fmt.Sprintf("New Balance: %v", acct.Balance)
	return "", errors.New(resp)
}

func resetPlay(text []string, decodedMap url.Values) (string, error) {
	defer measureTime(time.Now(), "reset")

	if len(text) != 1 { // Only reset accepted
		return "", errors.New("Error: Invalid number arguments")
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
	return "", fmt.Errorf("New Balance: %v", newAcct.Balance)
}

func portfolioPlay(text []string, decodedMap url.Values) (string, error) {
	defer measureTime(time.Now(), "portfolio")

	if len(text) != 1 { // Only portfolio accepted
		return "", errors.New("Error: Invalid number arguments")
	}

	// User and token to be used as lookup
	user := decodedMap["user_id"]
	token := decodedMap["token"]
	key := fmt.Sprintf("%v%v%v", "ACCT", token, user)

	client := connectRedis()

	acct, err := loadAccount(client, key)
	if err != nil {
		return "", errors.New("Unable to load account")
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
			return "", errors.New("Unable to get quotes")
		}

		rows := make([][]interface{}, 0, len(acct.Holdings))
		for _, quote := range info {
			h := acct.Holdings[quote.Ticker]
			total += float64(h.Shares) * quote.LatestPrice
			delta := float64(h.Shares) * (quote.LatestPrice - h.Strike)
			totalChange += delta
			deltaStr := fmt.Sprintf("%0.2f", delta)
			rows = append(rows, []interface{}{quote.Ticker, h.Shares, h.Strike, quote.LatestPrice, deltaStr})
		}

		rows = append(rows, []interface{}{"Total", "---", "---", "---", fmt.Sprintf("%0.2f", totalChange)})

		t := gotabulate.Create(rows)
		t.SetHeaders([]string{"Ticker", "Shares", "Strike", "Current", "Gain/Loss $"})
		t.SetAlign("left")
		t.SetHideLines([]string{"bottomLine", "betweenLine", "top"})
		table := t.Render("simple")
		summary := fmt.Sprintf("Portfolio Value: $%0.2f\nBalance: $%0.2f\nTotal: $%0.2f", total, acct.Balance, total+acct.Balance)
		resp := fmt.Sprintf("```%v\n%v```", table, summary)
		return resp, nil
	}

	s := fmt.Sprintf("Balance: $%0.2f", acct.Balance)
	return s, nil
}

func buyPlay(text []string, decodedMap url.Values) (string, error) {
	defer measureTime(time.Now(), "buy")

	if len(text) != 3 { // BUY TICKER SHARES
		return "", errors.New("Error: Invalid number arguments")
	}

	// Chop off buy arg
	text = text[1:]

	// Parse number of shares to purchase
	amt, err := strconv.ParseUint(text[1], 10, 64)
	if err != nil {
		return "", fmt.Errorf("Invalid amount: %v", err)
	}

	// lookup ticker price
	ticker := text[0]
	price, err := stockInterface.Price(ticker)
	if err != nil {
		return "", fmt.Errorf("Unable to get price: %v", err)
	}

	// User and token to be used as lookup
	user := decodedMap["user_id"]
	token := decodedMap["token"]
	key := fmt.Sprintf("%v%v%v", "ACCT", token, user)

	client := connectRedis()

	acct, err := loadAccount(client, key)
	if err != nil {
		return "", errors.New("Unable to load account")
	}

	// check if enough money in account
	if acct.Balance < (price * float64(amt)) {
		return "", errors.New("Insufficient funds")
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
		return "", fmt.Errorf("Unable to save account: %v", err)
	}

	return "", errors.New("Done")
}

func sellPlay(text []string, decodedMap url.Values) (string, error) {
	defer measureTime(time.Now(), "sell")

	if len(text) != 3 { // SELL TICKER SHARES
		return "", errors.New("Error: Invalid number arguments")
	}

	// Chop off buy arg
	text = text[1:]

	// Parse number of shares to purchase
	amt, err := strconv.ParseUint(text[1], 10, 64)
	if err != nil {
		return "", fmt.Errorf("Invalid amount: %v", err)
	}

	// lookup ticker price
	ticker := text[0]
	price, err := stockInterface.Price(ticker)
	if err != nil {
		return "", fmt.Errorf("Unable to get price: %v", err)
	}

	// User and token to be used as lookup
	user := decodedMap["user_id"]
	token := decodedMap["token"]
	key := fmt.Sprintf("%v%v%v", "ACCT", token, user)

	client := connectRedis()

	acct, err := loadAccount(client, key)
	if err != nil {
		return "", errors.New("Unable to load account")
	}

	h, ok := acct.Holdings[ticker]
	if !ok || h.Shares < amt {
		return "", errors.New("Not enough shares")
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
		return "", fmt.Errorf("Unable to save account: %v", err)
	}

	return "", errors.New("Done")
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
func getNews(text []string, decodedMap url.Values) (string, error) {
	defer measureTime(time.Now(), "news")

	if len(text) != 2 {
		return "", errors.New("Error: Invalid number of arguments")
	}
	// Chop off news arg
	text = text[1:]

	latestNews, err := stockInterface.News(text[0])
	if err != nil {
		return "", errors.New("Error: No news is good news right?")
	}

	// Try and pretty print them
	var printNews string
	for _, n := range latestNews {
		printNews = fmt.Sprintf("%s%s\n\n", printNews, n)
	}

	return printNews, nil
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

func getStats(text []string, _ url.Values) (string, error) {
	defer measureTime(time.Now(), "stats")

	// chop off stats arg
	text = text[1:]

	stats, err := stockInterface.Stats(text[0])
	if err != nil {
		return "", err
	}

	if len(text) == 1 { // user didn't request specific stats, return all of them
		return fmt.Sprintf("%v", stats), nil
	}

	// Pull out only the requested info
	requested := make(map[string]bool, len(text)-1)
	for i := 1; i < len(text); i++ {
		requested[strings.ToLower(text[i])] = true
	}

	var retStr string

	// Find the stat inside the struct
	v := reflect.ValueOf(stats).Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		n := strings.ToLower(v.Type().Field(i).Name)

		if requested[n] {
			retStr = retStr + fmt.Sprintf("%s: %v\n", n, f)
		}
	}

	return retStr, nil
}

// sortableList is a sort wrapper around a slice of stock quotes
// they are sorted by percent change
type sortableList []*stock.Quote

func (s sortableList) Len() int { return len(s) }

func (s sortableList) Less(i, j int) bool { return s[i].ChangePercent > s[j].ChangePercent }

func (s sortableList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
