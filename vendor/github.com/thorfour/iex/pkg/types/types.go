package types

import (
	"fmt"
	"strings"
)

const (
	APIVersion = "v1"
	QuoteStr   = "quote"
	NewsStr    = "news"
	StatsStr   = "stats"
	CompanyStr = "company"
	ChartStr   = "chart"
	StockStr   = "stock"
	PriceStr   = "price"
	BatchStr   = "batch"
	LastStr    = "last"
	MrktStr    = "market"
	APIURL     = "cloud.iexapis.com"
)

// Quote repesents the format returned for a quote from IEX(https://iextrading.com)
type Quote struct {
	Symbol                string  `json:"symbol"`
	CompanyName           string  `json:"companyName"`
	CalculationPrice      string  `json:"calculationPrice"`
	Open                  float64 `json:"open"`
	OpenTime              int64   `json:"openTime"`
	Close                 float64 `json:"close"`
	CloseTime             int64   `json:"closeTime"`
	High                  float64 `json:"high"`
	Low                   float64 `json:"low"`
	LatestPrice           float64 `json:"latestPrice"`
	LatestSource          string  `json:"latestSource"`
	LatestTime            string  `json:"latestTime"`
	LatestUpdate          int64   `json:"latestUpdate"`
	LatestVolume          int     `json:"latestVolume"`
	IexRealtimePrice      float64 `json:"iexRealtimePrice"`
	IexRealtimeSize       int     `json:"iexRealtimeSize"`
	IexLastUpdated        int64   `json:"iexLastUpdated"`
	DelayedPrice          float64 `json:"delayedPrice"`
	DelayedPriceTime      int64   `json:"delayedPriceTime"`
	ExtendedPrice         float64 `json:"extendedPrice"`
	ExtendedChange        float64 `json:"extendedChange"`
	ExtendedChangePercent float64 `json:"extendedChangePercent"`
	ExtendedPriceTime     int64   `json:"extendedPriceTime"`
	PreviousClose         float64 `json:"previousClose"`
	Change                float64 `json:"change"`
	ChangePercent         float64 `json:"changePercent"`
	IexMarketPercent      float64 `json:"iexMarketPercent"`
	IexVolume             int     `json:"iexVolume"`
	AvgTotalVolume        int     `json:"avgTotalVolume"`
	IexBidPrice           int     `json:"iexBidPrice"`
	IexBidSize            int     `json:"iexBidSize"`
	IexAskPrice           int     `json:"iexAskPrice"`
	IexAskSize            int     `json:"iexAskSize"`
	MarketCap             int64   `json:"marketCap"`
	PeRatio               float64 `json:"peRatio"`
	Week52High            float64 `json:"week52High"`
	Week52Low             float64 `json:"week52Low"`
	YtdChange             float64 `json:"ytdChange"`
}

// News is the news structure returned from IEX
type News struct {
	Datetime   int64  `json:"datetime"`
	Headline   string `json:"headline"`
	Source     string `json:"source"`
	URL        string `json:"url"`
	Summary    string `json:"summary"`
	Related    string `json:"related"`
	Image      string `json:"image"`
	Lang       string `json:"lang"`
	HasPaywall bool   `json:"hasPaywall"`
}

// Batch is a []Quote
type Batch map[string]map[string]Quote

// Quote returns the quote in a iex batch for a specific ticker
// returns error if symbol does not exist
func (i Batch) Quote(ticker string) (Quote, error) {
	ticker = strings.ToUpper(ticker)
	t, ok := i[ticker]
	if !ok {
		return Quote{}, fmt.Errorf("Failed to find %v in batch request", ticker)
	}

	return t["quote"], nil
}

// Stats holds information about a companys stats
type Stats struct {
	CompanyName         string      `json:"companyName"`
	Marketcap           int64       `json:"marketcap"`
	Beta                float64     `json:"beta"`
	Week52High          float64     `json:"week52high"`
	Week52Low           float64     `json:"week52low"`
	Week52Change        float64     `json:"week52change"`
	ShortInterest       float64     `json:"shortInterest"`
	ShortDate           interface{} `json:"shortDate"`
	DividendRate        float64     `json:"dividendRate"`
	DividendYield       float64     `json:"dividendYield"`
	ExDividendDate      interface{} `json:"exDividendDate"`
	LatestEPS           float64     `json:"latestEPS"`
	LatestEPSDate       interface{} `json:"latestEPSDate"`
	SharesOutstanding   float64     `json:"sharesOutstanding"`
	Float               float64     `json:"float"`
	ReturnOnEquity      float64     `json:"returnOnEquity"`
	ConsensusEPS        float64     `json:"consensusEPS"`
	NumberOfEstimates   float64     `json:"numberOfEstimates"`
	EPSSurpriseDollar   interface{} `json:"EPSSurpriseDollar"`
	EPSSurprisePercent  float64     `json:"EPSSurprisePercent"`
	Symbol              string      `json:"symbol"`
	EBITDA              float64     `json:"EBITDA"`
	Revenue             float64     `json:"revenue"`
	GrossProfit         float64     `json:"grossProfit"`
	Cash                float64     `json:"cash"`
	Debt                float64     `json:"debt"`
	TtmEPS              float64     `json:"ttmEPS"`
	RevenuePerShare     float64     `json:"revenuePerShare"`
	RevenuePerEmployee  float64     `json:"revenuePerEmployee"`
	PeRatioHigh         float64     `json:"peRatioHigh"`
	PeRatioLow          float64     `json:"peRatioLow"`
	ReturnOnAssets      float64     `json:"returnOnAssets"`
	ReturnOnCapital     interface{} `json:"returnOnCapital"`
	ProfitMargin        float64     `json:"profitMargin"`
	PriceToSales        float64     `json:"priceToSales"`
	PriceToBook         float64     `json:"priceToBook"`
	Day200MovingAvg     float64     `json:"day200MovingAvg"`
	Day50MovingAvg      float64     `json:"day50MovingAvg"`
	InstitutionPercent  float64     `json:"institutionPercent"`
	InsiderPercent      interface{} `json:"insiderPercent"`
	ShortRatio          interface{} `json:"shortRatio"`
	Year5ChangePercent  float64     `json:"year5ChangePercent"`
	Year2ChangePercent  float64     `json:"year2ChangePercent"`
	Year1ChangePercent  float64     `json:"year1ChangePercent"`
	YtdChangePercent    float64     `json:"ytdChangePercent"`
	Month6ChangePercent float64     `json:"month6ChangePercent"`
	Month3ChangePercent float64     `json:"month3ChangePercent"`
	Month1ChangePercent float64     `json:"month1ChangePercent"`
	Day5ChangePercent   float64     `json:"day5ChangePercent"`
	Day30ChangePercent  float64     `json:"day30ChangePercent"`
}

// Company contains company information
type Company struct {
	Symbol      string   `json:"symbol"`
	CompanyName string   `json:"companyName"`
	Exchange    string   `json:"exchange"`
	Industry    string   `json:"industry"`
	Website     string   `json:"website"`
	Description string   `json:"description"`
	CEO         string   `json:"CEO"`
	IssueType   string   `json:"issueType"`
	Sector      string   `json:"sector"`
	Tags        []string `json:"tags"`
}
