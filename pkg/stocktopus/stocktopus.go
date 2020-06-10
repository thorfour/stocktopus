package stocktopus

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/thorfour/stocktopus/pkg/cfg"
	"github.com/thorfour/stocktopus/pkg/stock"
	"gopkg.in/redis.v5"
)

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
	cmdHist = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "command_timings",
		Help: "A histogram of cmd request execution times",
	}, []string{"command"})
)

// measureTime is a helper function to measure the execution time of a function
func measureTime(start time.Time, label string) {
	cmdHist.WithLabelValues(label).Observe(time.Since(start).Seconds())
}

// Process url string to provide stocktpus functionality
func Process(args url.Values) (string, error) {
	s := &Stocktopus{ // TODO lift this out
		kvstore: redis.NewClient(&redis.Options{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPw,
			DB:       0,
		}),
		stockInterface: &stock.IexWrapper{},
	}
	text, ok := args["text"]
	if !ok {
		return "", errors.New("Bad request")
	}

	text = strings.Split(strings.ToUpper(text[0]), " ")
	switch text[0] {
	case buy:
		shares, err := strconv.Atoi(text[2])
		if err != nil {
			return "", err
		}
		s.Buy(text[1], uint64(shares), acctKey(args))
	case sell:
		shares, err := strconv.Atoi(text[2])
		if err != nil {
			return "", err
		}
		s.Sell(text[1], uint64(shares), acctKey(args))
	case deposit:
		amount, err := strconv.Atoi(text[0])
		if err != nil {
			return "", err
		}
		s.Deposit(float64(amount), acctKey(args))
	case portfolio:
		acct, err := s.Portfolio(acctKey(args))
		if err != nil {
			return "", err
		}

		return acct.String(), err
	case reset:
		return "", s.Clear(acctKey(args))
	default:
		wl, err := s.getQuotes(text)
		if err != nil {
			return "", err
		}

		return wl.String(), nil
	}

	return "", fmt.Errorf("bad request")
}

func listkey(text []string, decodedMap url.Values) string {

	// User and token to be used as watch list lookup
	user := decodedMap["user_id"]
	token := decodedMap["token"]

	// If the first arg starts with '#' then it's the name of the list
	if strings.HasPrefix(text[0], "#") {
		user = []string{strings.ToLower(text[0][1:]), decodedMap["team_id"][0]}
	}

	return fmt.Sprintf("%v%v", token, user)
}

func acctKey(decodedMap url.Values) string {
	// User and token to be used as lookup
	user := decodedMap["user_id"]
	token := decodedMap["token"]
	return fmt.Sprintf("%v%v%v", "ACCT", token, user)
}
