package stocktopus

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/alicebob/miniredis"
	redis "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
	"github.com/thorfour/iex/pkg/types"
	"github.com/thorfour/stocktopus/pkg/cfg"
	"github.com/thorfour/stocktopus/pkg/stock"
)

// fakeLookup implementes the stock.Lookup interface
type fakeLookup struct {
	fakeQuotes  []*stock.Quote
	fakeCompany *types.Company
	fakeStats   *types.Stats
	fakeNews    []string
}

func (f *fakeLookup) Price(string) (float64, error)                { return 1.00, nil }
func (f *fakeLookup) BatchQuotes([]string) ([]*stock.Quote, error) { return f.fakeQuotes, nil }
func (f *fakeLookup) News(string) ([]string, error)                { return f.fakeNews, nil }
func (f *fakeLookup) Stats(string) (*types.Stats, error)           { return f.fakeStats, nil }
func (f *fakeLookup) Company(string) (*types.Company, error)       { return f.fakeCompany, nil }

func TestCommands(t *testing.T) {

	// Start mini redis instance to connect to
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()
	cfg.RedisAddr = mr.Addr()

	s := &Stocktopus{
		kvstore: redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		}),
		stockInterface: &fakeLookup{
			fakeQuotes: []*stock.Quote{
				{
					Ticker:        "AMD",
					LatestPrice:   1.00,
					Change:        0,
					ChangePercent: 0,
				},
			},
			fakeCompany: &types.Company{},
			fakeStats:   &types.Stats{},
			fakeNews:    []string{},
		},
	}

	tests := []struct {
		name string
		text string
		err  error
	}{
		{
			name: "single quote",
			text: "amd",
			err:  nil,
		},
		{
			name: "retrieve empty group list",
			text: "list #mylist",
			err:  errors.New("No List"),
		},
		{
			name: "retrieve empty list",
			text: "list",
			err:  errors.New("No List"),
		},
		{
			name: "add to group list",
			text: "watch #mylist amd",
		},
		{
			name: "add to list",
			text: "watch amd",
		},
		{
			name: "retrieve group list",
			text: "list #mylist",
		},
		{
			name: "retrieve list",
			text: "list",
		},
		{
			name: "unwatch list",
			text: "unwatch amd",
		},
		{
			name: "unwatch group list",
			text: "unwatch #mylist amd",
		},
		{
			name: "clear list",
			text: "clear",
		},
		{
			name: "clear group list",
			text: "clear #mylist",
		},
		{
			name: "deposit",
			text: "deposit 100",
		},
		{
			name: "reset",
			text: "reset",
		},
		{
			name: "portfolio",
			text: "portfolio",
		},
		{
			name: "buy insufficient",
			text: "buy amd 1",
			err:  errors.New("Insufficient funds"),
		},
		{
			name: "deposit 1k",
			text: "deposit 1000",
		},
		{
			name: "buy amd",
			text: "buy amd 1",
		},
		{
			name: "portfolio with holdings",
			text: "portfolio",
		},
		{
			name: "sell amd too many",
			text: "sell amd 10",
			err:  errors.New("Not enough shares"),
		},
		{
			name: "sell amd",
			text: "sell amd 1",
		},
		{
			name: "info",
			text: "info amd",
		},
		{
			name: "stats",
			text: "stats amd",
		},
		{
			name: "stats with filter",
			text: "stats amd beta",
		},
		{
			name: "news",
			text: "news amd",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := url.Values{}
			v.Add("user_id", "test")
			v.Add("token", "token")
			v.Add("team_id", "team")
			v.Add("text", test.text)
			_, err := s.Process(v)
			require.Equal(t, test.err, err)
		})
	}
}

func TestAccount(t *testing.T) {

	// Start mini redis instance to connect to
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()
	cfg.RedisAddr = mr.Addr()

	s := &Stocktopus{
		kvstore: redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		}),
		stockInterface: &fakeLookup{
			fakeQuotes: []*stock.Quote{
				{
					Ticker:        "AMD",
					LatestPrice:   1.00,
					Change:        0,
					ChangePercent: 0,
				},
			},
			fakeCompany: &types.Company{},
			fakeStats:   &types.Stats{},
			fakeNews:    []string{},
		},
	}

	ctx := context.Background()
	a, err := s.Deposit(ctx, 1000, "mykey")
	require.NoError(t, err)
	require.Equal(t, &Account{
		Balance:  1000,
		Holdings: map[string]Holding{},
	}, a)

	require.Equal(t, "Balance: $1000.00", a.String())

	a, err = s.Buy(ctx, "AMD", 1, "mykey")
	require.NoError(t, err)
	require.Equal(t, &Account{
		Balance: 999,
		Holdings: map[string]Holding{
			"AMD": {
				Strike: 1,
				Shares: 1,
			},
		},
	}, a)

	a, err = s.Latest(ctx, a)
	require.NoError(t, err)
	require.Equal(t, &Account{
		Balance: 999,
		Holdings: map[string]Holding{
			"AMD": {
				Strike: 1,
				Shares: 1,
			},
		},
		Latest: map[string]float64{
			"AMD": 1.00,
		},
	}, a)

	exp :=
		` Ticker       Shares       Strike       Current       Gain/Loss $    
-----------  -----------  -----------  ------------  ----------------
 AMD          1            1            1             0.00           
 Total        ---          ---          ---           0.00           

Portfolio Value: $1.00
Balance: $999.00
Total: $1000.00`

	require.Equal(t, exp, a.String())
}
