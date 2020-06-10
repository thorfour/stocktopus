package stocktopus

import (
	"encoding/json"
	"errors"
	"fmt"

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

// Print returns a string representation of a watchlist
// TODO should it just return a struct that is a watchlist?
func (s *Stocktopus) Print(key string) (string, error) {
	list, err := s.kvstore.SMembers(key).Result()
	if err != nil {
		return "", fmt.Errorf("SMembers failed: %w", err)
	}

	if len(list) == 0 {
		return "", fmt.Errorf("No List")
	}

	return getQuotes(list)
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
func (s *Stocktopus) Deposit(amount float64, key string) (*account, error) {
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
		acct.Holdings[ticker] = holding{price, shares}
	} else {
		newShares := h.Shares + shares
		acct.Holdings[ticker] = holding{price, newShares}
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
		acct.Holdings[ticker] = holding{h.Strike, newShares}
	}

	acct.Balance += float64(shares) * price

	if err := s.saveAccount(key, acct); err != nil {
		return fmt.Errorf("Unable to save account: %w", err)
	}

	return nil
}

// Portfolio returns the account for a given key
func (s *Stocktopus) Portfolio(key string) (*account, error) {
	return s.account(key)
}

func (s *Stocktopus) account(key string) (*account, error) {
	serialized, err := s.kvstore.Get(key).Result()
	if err != nil {
		return nil, fmt.Errorf("Unable to load account: %w", err)
	}

	// Deserialize into struct
	acct := &account{}
	if err := json.Unmarshal([]byte(serialized), acct); err != nil {
		return nil, fmt.Errorf("Unable to parse account: %w", err)
	}

	return acct, nil
}

func (s *Stocktopus) saveAccount(key string, acct *account) error {
	b, err := json.Marshal(acct)
	if err != nil {
		return fmt.Errorf("Failed to serialize account: %w", err)
	}

	if _, err := s.kvstore.Set(key, b, 0).Result(); err != nil {
		return fmt.Errorf("Failed to save account: %w", err)
	}

	return nil
}
