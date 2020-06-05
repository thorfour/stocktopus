package stocktopus

import (
	"errors"
	"net/url"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/require"
	"github.com/thorfour/iex/pkg/types"
	"github.com/thorfour/stocktopus/pkg/cfg"
	"github.com/thorfour/stocktopus/pkg/stock"
)

// fakeLookup implementes the stock.Lookup interface
type fakeLookup struct {
	fakeQuotes []*stock.Quote
}

func (f *fakeLookup) Price(string) (float64, error)                { return 0, nil }
func (f *fakeLookup) BatchQuotes([]string) ([]*stock.Quote, error) { return f.fakeQuotes, nil }
func (f *fakeLookup) News(string) ([]string, error)                { return nil, nil }
func (f *fakeLookup) Stats(string) (*types.Stats, error)           { return nil, nil }
func (f *fakeLookup) Company(string) (*types.Company, error)       { return nil, nil }

func TestCommands(t *testing.T) {

	// Start mini redis instance to connect to
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()
	cfg.RedisAddr = mr.Addr()

	stockInterface = &fakeLookup{
		fakeQuotes: []*stock.Quote{
			{
				Ticker:        "AMD",
				LatestPrice:   1.00,
				Change:        0,
				ChangePercent: 0,
			},
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
			err:  errors.New("Error: No List"),
		},
		{
			name: "retrieve empty list",
			text: "list",
			err:  errors.New("Error: No List"),
		},
		{
			name: "add to group list",
			text: "watch #mylist amd",
			err:  errors.New("Added"),
		},
		{
			name: "add to list",
			text: "watch amd",
			err:  errors.New("Added"),
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
			err:  errors.New("Removed"),
		},
		{
			name: "unwatch group list",
			text: "unwatch #mylist amd",
			err:  errors.New("Removed"),
		},
		{
			name: "clear list",
			text: "clear",
			err:  errors.New("Removed"),
		},
		{
			name: "clear group list",
			text: "clear #mylist",
			err:  errors.New("Removed"),
		},
		{
			name: "deposit",
			text: "deposit 100",
			err:  errors.New("New Balance: 100"),
		},
		{
			name: "reset",
			text: "reset",
			err:  errors.New("New Balance: 0"),
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
			err:  errors.New("New Balance: 1000"),
		},
		{
			name: "buy amd",
			text: "buy amd 1",
			err:  errors.New("Done"),
		},
		{
			name: "sell amd too many",
			text: "sell amd 10",
			err:  errors.New("Not enough shares"),
		},
		{
			name: "sell amd",
			text: "sell amd 1",
			err:  errors.New("Done"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := url.Values{}
			v.Add("user_id", "test")
			v.Add("token", "token")
			v.Add("team_id", "team")
			v.Add("text", test.text)
			_, err := Process(v)
			require.Equal(t, test.err, err)
		})
	}
}
