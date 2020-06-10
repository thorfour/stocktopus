package stocktopus

import (
	"encoding/json"
	"errors"
	"fmt"

	redis "gopkg.in/redis.v5"
)

var (
	// ErrInvalidArguments is returned when the wrong number of args are passed in
	ErrInvalidArguments = errors.New("Error: invalid number of arguments")
)

// Stocktopus facilitates the retrieval and storage of stock information
type Stocktopus struct {
	kvstore *redis.Client
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
		return fmt.Errorf("SRem failed: %w", err)
	}

	return errors.New("Removed")
}

//-------------------------------------
//
// Play money actions
//
//-------------------------------------

// Deposit play money in account
func (s *Stocktopus) Deposit(amount float64, key string) (string, error) {
	acct, err := s.account(key)
	if err != nil {
		// TODO check for no key
		return "", err
	}
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
