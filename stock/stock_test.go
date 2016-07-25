package stock

import (
	"fmt"
	"testing"
)

func TestQuoteGoogle(t *testing.T) {

	resp, err := GetQuoteGoogle("AMD")
	if err != nil {
		t.Fail()
	}
	fmt.Println(resp)
	resp, err = GetQuoteGoogle("TWLO")
	if err != nil {
		t.Fail()
	}
	fmt.Println(resp)
	resp, err = GetQuoteGoogle("WDC")
	if err != nil {
		t.Fail()
	}
	fmt.Println(resp)
}

func TestQuoteMOD(t *testing.T) {

	resp, err := GetQuoteMOD("AMD")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(resp)
}
