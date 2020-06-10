package stocktopus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	redis "github.com/go-redis/redis/v8"
	"github.com/thorfour/iex/pkg/types"
	"github.com/thorfour/stocktopus/pkg/stock"
)

var (
	// ErrInvalidArguments is returned when the wrong number of args are passed in
	ErrInvalidArguments = errors.New("Error: invalid number of arguments")
)

// Stocktopus facilitates the retrieval and storage of stock information
type Stocktopus struct {
	KVStore        *redis.Client
	StockInterface stock.Lookup
}

//-------------------------------------
//
// Watch list commands
//
//-------------------------------------

// Add ticker(s) to a watch list
func (s *Stocktopus) Add(ctx context.Context, tickers []string, key string) error {
	if len(tickers) == 0 {
		return ErrInvalidArguments
	}

	if _, err := s.KVStore.SAdd(ctx, key, tickers).Result(); err != nil {
		return fmt.Errorf("SAdd failed: %w", err)
	}

	return nil
}

// Print returns a watchlist
func (s *Stocktopus) Print(ctx context.Context, key string) (WatchList, error) {
	list, err := s.KVStore.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("SMembers failed: %w", err)
	}

	if len(list) == 0 {
		return nil, fmt.Errorf("No List")
	}

	return s.GetQuotes(list)
}

// Remove ticker(s) from a watch list
func (s *Stocktopus) Remove(ctx context.Context, tickers []string, key string) error {
	if len(tickers) == 0 {
		return ErrInvalidArguments
	}

	if _, err := s.KVStore.SRem(ctx, key, tickers).Result(); err != nil {
		return fmt.Errorf("SRem failed: %w", err)
	}

	return nil
}

// Clear a watchlist or account by deleting the key from the KVStore
func (s *Stocktopus) Clear(ctx context.Context, key string) error {
	if _, err := s.KVStore.Del(ctx, key).Result(); err != nil {
		return fmt.Errorf("Del failed: %w", err)
	}

	return nil
}

//-------------------------------------
//
// Play money actions
//
//-------------------------------------

// Deposit play money in account
// NOTE: amount is a float because of legacy mistakes
func (s *Stocktopus) Deposit(ctx context.Context, amount float64, key string) (*Account, error) {
	acct, err := s.account(ctx, key)
	if err != nil {
		// TODO check for no key, because we want to open an account if there isn't one
		return nil, err
	}

	acct.Balance += amount

	if err := s.saveAccount(ctx, key, acct); err != nil {
		return nil, err
	}

	return acct, nil
}

// Buy shares for play money portfolio
func (s *Stocktopus) Buy(ctx context.Context, ticker string, shares uint64, key string) (*Account, error) {

	price, err := s.StockInterface.Price(ticker)
	if err != nil {
		return nil, fmt.Errorf("quote failed: %w", err)
	}

	acct, err := s.account(ctx, key)
	if err != nil {
		return nil, err
	}

	if acct.Balance < (price * float64(shares)) {
		return nil, errors.New("Insufficient funds")
	}

	// Add to account
	acct.Balance -= (price * float64(shares))
	h := acct.Holdings[ticker]
	acct.Holdings[ticker] = Holding{
		Strike: price,
		Shares: h.Shares + shares,
	}

	if err := s.saveAccount(ctx, key, acct); err != nil {
		return nil, err
	}

	return acct, nil
}

// Sell shares for play money portfolio
func (s *Stocktopus) Sell(ctx context.Context, ticker string, shares uint64, key string) (*Account, error) {

	price, err := s.StockInterface.Price(ticker)
	if err != nil {
		return nil, fmt.Errorf("quote failed: %w", err)
	}

	acct, err := s.account(ctx, key)
	if err != nil {
		return nil, err
	}

	h, ok := acct.Holdings[ticker]
	if !ok || h.Shares < shares {
		return nil, errors.New("Not enough shares")
	}

	newShares := h.Shares - shares
	if newShares == 0 {
		delete(acct.Holdings, ticker)
	} else {
		acct.Holdings[ticker] = Holding{
			Strike: h.Strike,
			Shares: newShares,
		}
	}

	acct.Balance += float64(shares) * price

	if err := s.saveAccount(ctx, key, acct); err != nil {
		return nil, fmt.Errorf("Unable to save account: %w", err)
	}

	return acct, nil
}

// Portfolio returns the account for a given key
func (s *Stocktopus) Portfolio(ctx context.Context, key string) (*Account, error) {
	return s.account(ctx, key)
}

// Latest populates the Latest map in the account (it is not saved)
func (s *Stocktopus) Latest(ctx context.Context, acct *Account) (*Account, error) {
	tickers := make([]string, 0, len(acct.Holdings))
	for ticker := range acct.Holdings {
		tickers = append(tickers, ticker)
	}
	quotes, err := s.GetQuotes(tickers)
	if err != nil {
		return nil, err
	}

	// Populate latest prices
	acct.Latest = map[string]float64{}
	for _, q := range quotes {
		acct.Latest[q.Ticker] = q.LatestPrice
	}

	return acct, nil
}

//-------------------------------------
//
// Info API
//
//-------------------------------------

// Info returns info about a given company
func (s *Stocktopus) Info(ticker string) (*types.Company, error) {
	info, err := s.StockInterface.Company(ticker)
	if err != nil {
		return nil, fmt.Errorf("Failed to get company info: %w", err)
	}

	return info, nil
}

// News returns the headlines for a given company
func (s *Stocktopus) News(ticker string) ([]string, error) {
	news, err := s.StockInterface.News(ticker)
	if err != nil {
		return nil, fmt.Errorf("Failed to get news: %w", err)
	}

	return news, nil
}

// Stats returns company statistics
func (s *Stocktopus) Stats(ticker string) (*types.Stats, error) {
	stats, err := s.StockInterface.Stats(ticker)
	if err != nil {
		return nil, fmt.Errorf("Failed to get stats: %w", err)
	}

	return stats, nil
}

//-------------------------------------
//
// Helper funtions
//
//-------------------------------------

func (s *Stocktopus) account(ctx context.Context, key string) (*Account, error) {
	serialized, err := s.KVStore.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) { // If no key found return fresh account
			return &Account{
				Holdings: map[string]Holding{},
			}, nil
		}
		return nil, fmt.Errorf("Unable to load account: %w", err)
	}

	// Deserialize into struct
	acct := &Account{}
	if err := json.Unmarshal([]byte(serialized), acct); err != nil {
		return nil, fmt.Errorf("Unable to parse account: %w", err)
	}

	return acct, nil
}

func (s *Stocktopus) saveAccount(ctx context.Context, key string, acct *Account) error {
	b, err := json.Marshal(acct)
	if err != nil {
		return fmt.Errorf("Failed to serialize account: %w", err)
	}

	if _, err := s.KVStore.Set(ctx, key, b, 0).Result(); err != nil {
		return fmt.Errorf("Failed to save account: %w", err)
	}

	return nil
}

// GetQuotes returns a list of quotes from tickers
func (s *Stocktopus) GetQuotes(tickers []string) (WatchList, error) {
	quotes, err := s.StockInterface.BatchQuotes(tickers)
	if err != nil {
		return nil, err
	}

	// Sort the list
	sort.Sort(WatchList(quotes))

	return WatchList(quotes), nil
}

func (s *Stocktopus) getChartLink(ticker string) string {
	symbol := strings.ToUpper(ticker)
	return fmt.Sprintf("http://finviz.com/chart.ashx?t=%s&ty=c&ta=1&p=d&s=l", symbol)
}
