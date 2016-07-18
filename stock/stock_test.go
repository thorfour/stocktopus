package stock

import (
	"fmt"
	"testing"
)

func TestQuote(t *testing.T) {

	resp := GetQuote("AMD")
	fmt.Println(resp)
	resp = GetQuote("TWLO")
	fmt.Println(resp)
	resp = GetQuote("WDC")
	fmt.Println(resp)
}
