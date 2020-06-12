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
	"github.com/sirupsen/logrus"
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

// Response is the json struct for a slack response
type Response struct {
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
		logrus.WithField("msg", "error parse form").Error(err)
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	msg, err := s.Process(ctx, req.Form)
	if err != nil {
		logrus.WithField("msg", "Process failure").Error(err)
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(resp).Encode(msg); err != nil {
		logrus.WithField("msg", "encoding failure").Error(err)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Process a slack request
func (s *SlashServer) Process(ctx context.Context, args url.Values) (*Response, error) {
	text, ok := args["text"]
	if !ok {
		return nil, errors.New("Bad request")
	}

	text = strings.Split(strings.ToUpper(text[0]), " ")
	return s.command(ctx, text[0], text[1:], args)
}

// Command processes a stocktopus command
func (s *SlashServer) command(ctx context.Context, cmd string, args []string, info map[string][]string) (*Response, error) {

	defer s.measureTime(time.Now(), cmd)

	switch cmd {
	case buy:
		shares, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, err
		}

		if _, err := s.s.Buy(ctx, args[0], uint64(shares), acctKey(info)); err != nil {
			return nil, fmt.Errorf("Buy failed: %w", err)
		}

		return &Response{
			ResponseType: ephemeral,
			Text:         "Done",
		}, nil

	case sell:
		shares, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, err
		}

		if _, err := s.s.Sell(ctx, args[0], uint64(shares), acctKey(info)); err != nil {
			return nil, fmt.Errorf("Sell failed: %w", err)
		}

		return &Response{
			ResponseType: ephemeral,
			Text:         "Done",
		}, nil

	case deposit:
		amount, err := strconv.Atoi(args[0])
		if err != nil {
			return nil, err
		}

		a, err := s.s.Deposit(ctx, float64(amount), acctKey(info))
		if err != nil {
			return nil, fmt.Errorf("Deposit failed: %w", err)
		}

		return &Response{
			ResponseType: ephemeral,
			Text:         fmt.Sprintf("New Balance: %v", a.Balance),
		}, nil

	case portfolio:
		a, err := s.s.Portfolio(ctx, acctKey(info))
		if err != nil {
			return nil, fmt.Errorf("Portfolio failed: %w", err)
		}

		a, err = s.s.Latest(ctx, a)
		if err != nil {
			return nil, fmt.Errorf("Latest failed: %w", err)
		}

		return &Response{
			ResponseType: inchannel,
			Text:         fmt.Sprintf("```%v```", a.String()),
		}, nil

	case reset:
		if err := s.s.Clear(ctx, acctKey(info)); err != nil {
			return nil, fmt.Errorf("Clear failed: %w", err)
		}

		return &Response{
			ResponseType: ephemeral,
			Text:         "New Balance: 0",
		}, nil

	case addToList:
		if err := s.s.Add(ctx, args, listkey(args, info)); err != nil {
			return nil, fmt.Errorf("Add failed: %w", err)
		}

		return &Response{
			ResponseType: ephemeral,
			Text:         "Added",
		}, nil

	case printList:
		a, err := s.s.Print(ctx, listkey(args, info))
		if err != nil {
			return nil, fmt.Errorf("Print failed: %w", err)
		}

		// TODO get chart link

		return &Response{
			ResponseType: inchannel,
			Text:         fmt.Sprintf("```%v```", a.String()),
		}, nil

	case removeFromList:
		if err := s.s.Remove(ctx, args, listkey(args, info)); err != nil {
			return nil, fmt.Errorf("Remove failed: %w", err)
		}

		return &Response{
			ResponseType: ephemeral,
			Text:         "Removed",
		}, nil
	case clear:
		if err := s.s.Clear(ctx, listkey(args, info)); err != nil {
			return nil, fmt.Errorf("Clear failed: %w", err)
		}

		return &Response{
			ResponseType: ephemeral,
			Text:         "Removed",
		}, nil

	case infoCmd:
		c, err := s.s.Info(args[0])
		if err != nil {
			return nil, fmt.Errorf("Info failed: %w", err)
		}

		return &Response{
			ResponseType: inchannel,
			Text:         strings.Join([]string{c.CompanyName, c.Industry, c.Website, c.CEO, c.Description}, "\n"),
		}, nil

	case news:
		news, err := s.s.News(args[0])
		if err != nil {
			return nil, fmt.Errorf("News failed: %w", err)
		}

		return &Response{
			ResponseType: inchannel,
			Text:         strings.Join(news, "\n\n"),
		}, nil

	case stats:
		stats, err := s.s.Stats(args[0])
		if err != nil {
			return nil, fmt.Errorf("Stats failed: %w", err)
		}

		// TODO filter stats?

		return &Response{
			ResponseType: inchannel,
			Text:         fmt.Sprintf("```%v```", stocktopus.Stats(stats)),
		}, nil

	default:
		wl, err := s.s.GetQuotes(args)
		if err != nil {
			return nil, fmt.Errorf("GetQuotes failed: %w", err)
		}

		return &Response{
			ResponseType: inchannel,
			Text:         fmt.Sprintf("```%v```", wl.String()),
		}, nil
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
