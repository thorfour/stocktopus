package stocktopus

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Supported commands
const (
	addToList      = "WATCH"
	printList      = "LIST"
	removeFromList = "UNWATCH"
	clear          = "CLEAR"
	help           = "HELP"
	infoCmd        = "INFO"
	news           = "NEWS"
	stats          = "STATS"

	// Play money commands
	buy       = "BUY"
	sell      = "SELL"
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

// Command processes a stocktopus command
// TODO this should get lifted out
func (s *Stocktopus) Command(ctx context.Context, cmd string, args []string, info map[string][]string) (string, error) {

	switch cmd {
	case buy:
		shares, err := strconv.Atoi(args[1])
		if err != nil {
			return "", err
		}
		a, err := s.Buy(ctx, args[0], uint64(shares), acctKey(info))
		if err != nil {
			return "", err
		}

		return a.String(), nil
	case sell:
		shares, err := strconv.Atoi(args[1])
		if err != nil {
			return "", err
		}
		a, err := s.Sell(ctx, args[0], uint64(shares), acctKey(info))
		if err != nil {
			return "", err
		}

		return a.String(), nil
	case deposit:
		amount, err := strconv.Atoi(args[0])
		if err != nil {
			return "", err
		}
		a, err := s.Deposit(ctx, float64(amount), acctKey(info))
		if err != nil {
			return "", err
		}

		return a.String(), nil
	case portfolio:
		a, err := s.Portfolio(ctx, acctKey(info))
		if err != nil {
			return "", err
		}

		a, err = s.Latest(ctx, a)
		if err != nil {
			return "", err
		}

		return a.String(), nil
	case reset:
		return "", s.Clear(ctx, acctKey(info))
	case addToList:
		return "", s.Add(ctx, args, listkey(args, info))
	case printList:
		acct, err := s.Print(ctx, listkey(args, info))
		if err != nil {
			return "", err
		}

		return acct.String(), nil
	case removeFromList:
		return "", s.Remove(ctx, args, listkey(args, info))
	case clear:
		return "", s.Clear(ctx, listkey(args, info))
	case infoCmd:
		c, err := s.Info(args[0])
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%v", c), nil // TODO
	case news:
		news, err := s.News(args[0])
		if err != nil {
			return "", err
		}

		return strings.Join(news, "\n"), nil
	case stats:
		stats, err := s.Stats(args[0])
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%v", stats), nil
	default:
		wl, err := s.getQuotes(args)
		if err != nil {
			return "", err
		}

		return wl.String(), nil
	}
}

func listkey(text []string, decodedMap url.Values) string {

	// User and token to be used as watch list lookup
	user := decodedMap["user_id"]
	token := decodedMap["token"]

	// If the first arg starts with '#' then it's the name of the list
	if len(text) > 0 && strings.HasPrefix(text[0], "#") {
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
