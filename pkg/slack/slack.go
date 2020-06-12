package slack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/thorfour/stocktopus/pkg/stock"
	"github.com/thorfour/stocktopus/pkg/stocktopus"
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

const (
	ephemeral = "ephemeral"
	inchannel = "in_channel"
)

// response is the json struct for a slack response
type response struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

// SlashServer is a slack server that handles slash commands
type SlashServer struct {
	s       *stocktopus.Stocktopus
	cmdHist *prometheus.HistogramVec
}

// measureTime is a helper function to measure the execution time of a function
func (s *SlashServer) measureTime(start time.Time, label string) {
	s.cmdHist.WithLabelValues(label).Observe(time.Since(start).Seconds())
}

// New returns a new slash server
func New(kvstore *redis.Client, stocks stock.Lookup) *SlashServer {
	return &SlashServer{
		s: &stocktopus.Stocktopus{
			KVStore:        kvstore,
			StockInterface: stocks,
		},
		cmdHist: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: "command_timings",
			Help: "A histogram of cmd request execution times",
		},
			[]string{"command"},
		),
	}
}

// Handler is a http handler func for processing slack slash requests for stocktopus
func (s *SlashServer) Handler(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	msg, err := s.Process(ctx, req.Form)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	s.newResponse(resp, msg, nil)
}

// Process a slack request
func (s *SlashServer) Process(ctx context.Context, args url.Values) (string, error) {
	text, ok := args["text"]
	if !ok {
		return "", errors.New("Bad request")
	}

	text = strings.Split(strings.ToUpper(text[0]), " ")
	return s.Command(ctx, text[0], text[1:], args)
}

// TODO determine ephermeralness of response
// TODO some of these need to be wrapped with ```
func (s *SlashServer) newResponse(resp http.ResponseWriter, message string, err error) {
	r := &response{
		ResponseType: inchannel,
		Text:         message,
	}

	// Switch to an ephemeral message
	if err != nil {
		r.ResponseType = ephemeral
		r.Text = err.Error()
	}

	b, err := json.Marshal(r)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(b)
	return
}

// Command processes a stocktopus command
func (s *SlashServer) Command(ctx context.Context, cmd string, args []string, info map[string][]string) (string, error) {

	defer s.measureTime(time.Now(), cmd)

	switch cmd {
	case buy:
		shares, err := strconv.Atoi(args[1])
		if err != nil {
			return "", err
		}
		a, err := s.s.Buy(ctx, args[0], uint64(shares), acctKey(info))
		if err != nil {
			return "", err
		}

		return a.String(), nil
	case sell:
		shares, err := strconv.Atoi(args[1])
		if err != nil {
			return "", err
		}
		a, err := s.s.Sell(ctx, args[0], uint64(shares), acctKey(info))
		if err != nil {
			return "", err
		}

		return a.String(), nil
	case deposit:
		amount, err := strconv.Atoi(args[0])
		if err != nil {
			return "", err
		}
		a, err := s.s.Deposit(ctx, float64(amount), acctKey(info))
		if err != nil {
			return "", err
		}

		return a.String(), nil
	case portfolio:
		a, err := s.s.Portfolio(ctx, acctKey(info))
		if err != nil {
			return "", err
		}

		a, err = s.s.Latest(ctx, a)
		if err != nil {
			return "", err
		}

		return a.String(), nil
	case reset:
		return "", s.s.Clear(ctx, acctKey(info))
	case addToList:
		return "", s.s.Add(ctx, args, listkey(args, info))
	case printList:
		acct, err := s.s.Print(ctx, listkey(args, info))
		if err != nil {
			return "", err
		}

		return acct.String(), nil
	case removeFromList:
		return "", s.s.Remove(ctx, args, listkey(args, info))
	case clear:
		return "", s.s.Clear(ctx, listkey(args, info))
	case infoCmd:
		c, err := s.s.Info(args[0])
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%v", c), nil // TODO
	case news:
		news, err := s.s.News(args[0])
		if err != nil {
			return "", err
		}

		return strings.Join(news, "\n"), nil
	case stats:
		stats, err := s.s.Stats(args[0])
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%v", stats), nil
	default:
		wl, err := s.s.GetQuotes(args)
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
