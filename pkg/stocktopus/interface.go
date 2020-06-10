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

// Account is a users play money account
type Account struct {
	Balance  float64
	Holdings map[string]Holding
	Latest   map[string]float64
}

// Holding is a specific stock holding
type Holding struct {
	Strike float64
	Shares uint64
}

func (a *Account) String() string {
	if len(a.Holdings) <= 0 {
		return fmt.Sprintf("Balance: $%0.2f", a.Balance)
	}

	total := float64(0)
	totalChange := float64(0)
	rows := make([][]interface{}, 0, len(a.Holdings))
	for ticker, h := range a.Holdings {
		latest, ok := a.Latest[ticker]
		if !ok {
			continue
		}
		total += float64(h.Shares) * latest
		delta := float64(h.Shares) * (latest - h.Strike)
		totalChange += delta
		deltaStr := fmt.Sprintf("%0.2f", delta)
		rows = append(rows,
			[]interface{}{
				ticker,
				h.Shares,
				h.Strike,
				latest,
				deltaStr,
			},
		)
	}

	rows = append(rows,
		[]interface{}{
			"Total",
			"---",
			"---",
			"---",
			fmt.Sprintf("%0.2f", totalChange),
		},
	)

	t := gotabulate.Create(rows)
	t.SetHeaders([]string{"Ticker", "Shares", "Strike", "Current", "Gain/Loss $"})
	t.SetAlign("left")
	t.SetHideLines([]string{"bottomLine", "betweenLine", "top"})
	table := t.Render("simple")
	summary := fmt.Sprintf("Portfolio Value: $%0.2f\nBalance: $%0.2f\nTotal: $%0.2f", total, a.Balance, total+a.Balance)
	return fmt.Sprintf("%v\n%v", table, summary)
}
