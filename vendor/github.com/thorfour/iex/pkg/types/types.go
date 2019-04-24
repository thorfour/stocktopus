package types

import (
	"fmt"
	"strings"
)

const (
	APIVersion = "1.0"
	QuoteStr   = "quote"
	NewsStr    = "news"
	StatsStr   = "stats"
	ChartStr   = "chart"
	StockStr   = "stock"
	PriceStr   = "price"
	BatchStr   = "batch"
	LastStr    = "last"
	MrktStr    = "market"
	APIURL     = "https://api.iextrading.com/"
)

// Quote repesents the format returned for a quote from IEX(https://iextrading.com)
type Quote struct {
	Symbol           string  `json:symbol`
	CompanyName      string  `json:companyName`
	PrimaryExchange  string  `json:primaryExchange`
	CalculationPrice string  `json:calculationPrice`
	IexRealtimePrice float64 `json:iexRealtimePrice`
	IexRealtimeSize  float64 `json:iexRealtimeSize`
	IexLastUpdated   float64 `json:iexLastUpdated`
	DelayedPrice     float64 `json:delayedPrice`
	DelayedPriceTime float64 `json:delayedPriceTime`
	PreviousClose    float64 `json:previousClose`
	Change           float64 `json:change`
	ChangePercent    float64 `json:changePercent`
	IexMarketPercent float64 `json:iexMarketPercent`
	IexVolume        float64 `json:iexVolume`
	AvgTotalVolume   float64 `json:avgTotalVolume`
	IexBidPrice      float64 `json:iexBidPrice`
	IexBidSize       float64 `json:iexBidSize`
	IexAskPrice      float64 `json:iexAskPrice`
	IexAskSize       float64 `json:iexAskSize`
	MarketCap        float64 `json:marketCap`
	LatestPrice      float64 `json:latestPrice`
	//PeRatio          float64 `json:peRatio`
	Week52High float64 `json:week52High`
	Week52Low  float64 `json:week52Low`
}

// News is the news structure returned from IEX
type News struct {
	DateTime string `json:datetime`
	Headline string `json:headline`
	Source   string `json:source`
	URL      string `json:url`
	Summary  string `json:summar`
	Related  string `json:related`
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

	q, ok := t[QuoteStr]
	if !ok {
		return Quote{}, fmt.Errorf("Failed to find quote for %v in batch request", ticker)
	}

	return q, nil
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
	ShortDate           float64     `json:"shortDate"`
	DividendRate        float64     `json:"dividendRate"`
	DividendYield       float64     `json:"dividendYield"`
	ExDividendDate      string      `json:"exDividendDate"`
	LatestEPS           float64     `json:"latestEPS"`
	LatestEPSDate       string      `json:"latestEPSDate"`
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
