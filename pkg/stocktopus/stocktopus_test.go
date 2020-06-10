package stocktopus

import (
	"context"
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

func TestAccount(t *testing.T) {

	// Start mini redis instance to connect to
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()
	cfg.RedisAddr = mr.Addr()

	s := &Stocktopus{
		KVStore: redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		}),
		StockInterface: &fakeLookup{
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

	a, err = s.Sell(ctx, "AMD", 1, "mykey")
	require.NoError(t, err)
	require.Equal(t, &Account{
		Balance:  1000,
		Holdings: map[string]Holding{},
	}, a)
	require.Equal(t, "Balance: $1000.00", a.String())
}
