package endpoint

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/thorfour/iex/pkg/types"
)

// Token used to make API requests
var Token string

func init() {
	Token = os.Getenv("IEX_API_TOKEN")
}

// API is the wrapper around a IEX URL
type API struct {
	u *url.URL
}

// Endpoint returns the url api endpint
func Endpoint() *API {
	return &API{&url.URL{
		Scheme: "https",
		Host:   types.APIURL,
		Path:   types.APIVersion,
	}}
}

// Market adds the market type
func (u API) Market() API {
	u.u.Path = path.Join(u.u.Path, types.MrktStr)
	return u
}

// Stock adds the stock type
func (u API) Stock() API {
	u.u.Path = path.Join(u.u.Path, types.StockStr)
	return u
}

// Quote adds the quote type
func (u API) Quote() API {
	u.u.Path = path.Join(u.u.Path, types.QuoteStr)
	return u
}

// Ticker adds the ticker
func (u API) Ticker(ticker string) API {
	u.u.Path = path.Join(u.u.Path, ticker)
	return u
}

// Price adds the price type
func (u API) Price() API {
	u.u.Path = path.Join(u.u.Path, types.PriceStr)
	return u
}

// Batch adds the batch type
func (u API) Batch() API {
	u.u.Path = path.Join(u.u.Path, types.BatchStr)
	return u
}

// Tickers adds the tickers as a comma separated list
func (u API) Tickers(t []string) API {
	q := u.u.Query()
	q.Add("symbols", strings.Join(t, ","))
	u.u.RawQuery = q.Encode()
	return u
}

// Types adds the types= argument
func (u API) Types(t ...string) API {
	q := u.u.Query()
	q.Add("types", strings.Join(t, ","))
	u.u.RawQuery = q.Encode()
	return u
}

// News adds the news type
func (u API) News() API {
	u.u.Path = path.Join(u.u.Path, types.NewsStr)
	return u
}

// Last adds the last type
func (u API) Last() API {
	u.u.Path = path.Join(u.u.Path, types.LastStr)
	return u
}

// Integer adds a integer argument
func (u API) Integer(a int) API {
	u.u.Path = path.Join(u.u.Path, fmt.Sprintf("%v", a))
	return u
}

// Stats adds the stats type
func (u API) Stats() API {
	u.u.Path = path.Join(u.u.Path, types.StatsStr)
	return u
}

// Company adds the company type
func (u API) Company() API {
	u.u.Path = path.Join(u.u.Path, types.CompanyStr)
	return u
}

func (u API) String() string {
	q := u.u.Query()
	q.Add("token", Token)
	u.u.RawQuery = q.Encode()
	return u.u.String()
}
