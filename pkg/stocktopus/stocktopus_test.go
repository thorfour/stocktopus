package stocktopus

import (
	"fmt"
	"testing"
)

func TestGetQuotesDelimited(t *testing.T) {
	_, err := getQuotes("brk.a", nil)
	if err != nil {
		t.Error(err)
	}
}

func TestGetQuotes(t *testing.T) {
	_, err := getQuotes("tsla amd wdc intc gpro f goog", nil)
	if err != nil {
		t.Error(err)
	}
}

func TestGetQuotesSingle(t *testing.T) {
	_, err := getQuotes("tsla", nil)
	if err != nil {
		t.Error(err)
	}
}

func TestGetQuotesSingleWithCurrency(t *testing.T) {
	_, err := getQuotes("tsla btcusd amd", nil)
	if err != nil {
		t.Error(err)
	}
}

func TestGetQuotesBad(t *testing.T) {
	if _, err := getQuotes("tsla osghoevcmi amd", nil); err != nil {
		t.Error(err)
	}
	if q, err := getQuotes("osghoevcmi", nil); err == nil {
		fmt.Println(q)
		t.Error("expected failure")
	}
	if _, err := getQuotes("osghoevcmi amd tsla", nil); err != nil {
		t.Error("expected failure")
	}
	if _, err := getQuotes("tsla amd aorghreqcm", nil); err != nil {
		t.Error("expected failure")
	}
	if _, err := getQuotes("amd aorghreqcm", nil); err != nil {
		t.Error("expected failure")
	}
}

func TestGetNews(t *testing.T) {
	if _, err := getNews([]string{"news", "amd"}, nil); err != nil {
		t.Error(err)
	}
}

func TestQuoteExchange(t *testing.T) {
	t.Skip("This test is not supported on IEX")
	_, err := getQuotes("TSE:ARE", nil)
	if err != nil {
		t.Error(err)
	}
}

func TestStats(t *testing.T) {
	if _, err := getStats([]string{"stats", "aapl", "marketcap"}, nil); err != nil {
		t.Error(err)
	}
}
