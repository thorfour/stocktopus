package stock

import (
	"fmt"
	"testing"
)

func TestQuote(t *testing.T) {

	resp, err := GetQuote("AMD")
	if err != nil {
		t.Fail()
	}
	fmt.Println(resp)
	resp, err = GetQuote("TWLO")
	if err != nil {
		t.Fail()
	}
	fmt.Println(resp)
	resp, err = GetQuote("WDC")
	if err != nil {
		t.Fail()
	}
	fmt.Println(resp)
}
