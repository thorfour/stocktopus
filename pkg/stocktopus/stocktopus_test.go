package stocktopus

import (
	"errors"
	"fmt"
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

	tests := map[string]struct {
		text string
		err  error
	}{
		"single quote": {
			text: "amd",
			err:  nil,
		},
		"add to group list": {
			text: "watch #mylist amd",
			err:  errors.New("Added"),
		},
		"add to list": {
			text: "watch amd",
			err:  errors.New("Added"),
		},
		"retrieve list": {
			text: "list #mylist",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			v := url.Values{}
			v.Add("user_id", "test")
			v.Add("token", "token")
			v.Add("team_id", "team")
			v.Add("text", test.text)
			_, err := Process(v)
			require.Equal(t, test.err, err)
		})
	}

	client := connectRedis()
	fmt.Println(client.Keys("*").String())
}
