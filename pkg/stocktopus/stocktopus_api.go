package stocktopus

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/thorfour/iex/pkg/types"
	"github.com/thorfour/stocktopus/pkg/stock"
	redis "gopkg.in/redis.v5"
)

var (
	// ErrInvalidArguments is returned when the wrong number of args are passed in
	ErrInvalidArguments = errors.New("Error: invalid number of arguments")
)

// Stocktopus facilitates the retrieval and storage of stock information
type Stocktopus struct {
	kvstore        *redis.Client
	stockInterface stock.Lookup
}

//-------------------------------------
//
// Watch list commands
//
//-------------------------------------

// Add ticker(s) to a watch list
func (s *Stocktopus) Add(tickers []string, key string) error {
	if len(tickers) == 0 {
		return ErrInvalidArguments
	}

	if _, err := s.kvstore.SAdd(key, tickers).Result(); err != nil {
		return fmt.Errorf("SAdd failed: %w", err)
	}

	return nil
}

// Print returns a watchlist
func (s *Stocktopus) Print(key string) (WatchList, error) {
	list, err := s.kvstore.SMembers(key).Result()
	if err != nil {
		return nil, fmt.Errorf("SMembers failed: %w", err)
	}

	if len(list) == 0 {
		return nil, fmt.Errorf("No List")
	}

	return s.getQuotes(list)
}

// Remove ticker(s) from a watch list
func (s *Stocktopus) Remove(tickers []string, key string) error {
	if len(tickers) == 0 {
		return ErrInvalidArguments
	}

	if _, err := s.kvstore.SRem(key, tickers).Result(); err != nil {
		return fmt.Errorf("SRem failed: %w", err)
	}

	return nil
}

// Clear a watchlist by deleting the key from the kvstore
func (s *Stocktopus) Clear(key string) error {
	if _, err := s.kvstore.Del(key).Result(); err != nil {
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
func (s *Stocktopus) Deposit(amount float64, key string) (*Account, error) {
	acct, err := s.account(key)
	if err != nil {
		// TODO check for no key, because we want to open an account if there isn't one
		return nil, err
	}

	acct.Balance += amount

	if err := s.saveAccount(key, acct); err != nil {
		return nil, err
	}

	return acct, nil
}

// Buy shares for play money portfolio
func (s *Stocktopus) Buy(ticker string, shares uint64, key string) error {

	price, err := s.stockInterface.Price(ticker)
	if err != nil {
		return fmt.Errorf("quote failed: %w", err)
	}

	acct, err := s.account(key)
	if err != nil {
		return err
	}

	if acct.Balance < (price * float64(shares)) {
		return errors.New("Insufficient funds")
	}

	// Add to account
	acct.Balance -= (price * float64(shares))
	h, ok := acct.Holdings[ticker]
	if !ok {
		acct.Holdings[ticker] = Holding{price, shares}
	} else {
		newShares := h.Shares + shares
		acct.Holdings[ticker] = Holding{price, newShares}
	}

	if err := s.saveAccount(key, acct); err != nil {
		return err
	}

	return nil
}

// Sell shares for play money portfolio
func (s *Stocktopus) Sell(ticker string, shares uint64, key string) error {

	price, err := s.stockInterface.Price(ticker)
	if err != nil {
		return fmt.Errorf("quote failed: %w", err)
	}

	acct, err := s.account(key)
	if err != nil {
		return err
	}

	h, ok := acct.Holdings[ticker]
	if !ok || h.Shares < shares {
		return errors.New("Not enough shares")
	}

	newShares := h.Shares - shares
	if newShares == 0 {
		delete(acct.Holdings, ticker)
	} else {
		acct.Holdings[ticker] = Holding{h.Strike, newShares}
	}

	acct.Balance += float64(shares) * price

	if err := s.saveAccount(key, acct); err != nil {
		return fmt.Errorf("Unable to save account: %w", err)
	}

	return nil
}

// Portfolio returns the account for a given key
func (s *Stocktopus) Portfolio(key string) (*Account, error) {
	return s.account(key)
}

//-------------------------------------
//
// Info API
//
//-------------------------------------

// Info returns info about a given company
func (s *Stocktopus) Info(ticker string) (*types.Company, error) {
	info, err := s.stockInterface.Company(ticker)
	if err != nil {
		return nil, fmt.Errorf("Failed to get company info: %w", err)
	}

	return info, nil
}

// News returns the headlines for a given company
func (s *Stocktopus) News(ticker string) ([]string, error) {
	news, err := s.stockInterface.News(ticker)
	if err != nil {
		return nil, fmt.Errorf("Failed to get news: %w", err)
	}

	return news, nil
}

// Stats returns company statistics
func (s *Stocktopus) Stats(ticker string) (*types.Stats, error) {
	stats, err := s.stockInterface.Stats(ticker)
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

func (s *Stocktopus) account(key string) (*Account, error) {
	serialized, err := s.kvstore.Get(key).Result()
	if err != nil {
		return nil, fmt.Errorf("Unable to load account: %w", err)
	}

	// Deserialize into struct
	acct := &Account{}
	if err := json.Unmarshal([]byte(serialized), acct); err != nil {
		return nil, fmt.Errorf("Unable to parse account: %w", err)
	}

	return acct, nil
}

func (s *Stocktopus) saveAccount(key string, acct *Account) error {
	b, err := json.Marshal(acct)
	if err != nil {
		return fmt.Errorf("Failed to serialize account: %w", err)
	}

	if _, err := s.kvstore.Set(key, b, 0).Result(); err != nil {
		return fmt.Errorf("Failed to save account: %w", err)
	}

	return nil
}

func (s *Stocktopus) getQuotes(tickers []string) (WatchList, error) {
	quotes, err := s.stockInterface.BatchQuotes(tickers)
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
