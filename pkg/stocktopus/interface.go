package stocktopus

import (
	"fmt"

	"github.com/bndr/gotabulate"
	"github.com/thorfour/stocktopus/pkg/stock"
)

// WatchListAPI is the interface for interacting with a watch list
type WatchListAPI interface {
	Add(tickers []string, key string) error
	Print(key string) (string, error)
	Remove(tickers []string, key string) error
	Clear(key string) error
}

// WatchList is a sort wrapper around a slice of stock quotes
// they are sorted by percent change
type WatchList []*stock.Quote

func (w WatchList) Len() int { return len(w) }

func (w WatchList) Less(i, j int) bool { return w[i].ChangePercent > w[j].ChangePercent }

func (w WatchList) Swap(i, j int) { w[i], w[j] = w[j], w[i] }

func (w WatchList) String() string {
	rows := make([][]interface{}, 0, len(w))
	cumsum := float64(0)
	for _, quote := range w {
		rows = append(rows,
			[]interface{}{
				quote.Ticker,
				quote.LatestPrice,
				fmt.Sprintf("%0.2f", quote.Change),
				fmt.Sprintf("%0.3f", (100 * quote.ChangePercent)),
			},
		)
		cumsum += (100 * quote.ChangePercent)
	}

	// Add an average row at the bottom
	rows = append(rows,
		[]interface{}{
			"Avg.",
			"---",
			"---",
			fmt.Sprintf("%0.3f%%", cumsum/float64(len(rows))),
		},
	)

	t := gotabulate.Create(rows)
	t.SetHeaders([]string{"Company", "Current Price", "Todays Change", "Percent Change"})
	t.SetAlign("right")
	t.SetHideLines([]string{"bottomLine", "betweenLine", "top"})

	return t.Render("simple")
}
