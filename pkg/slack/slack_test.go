package slack

import (
	"flag"
	"fmt"
	"testing"
)

var (
	token = flag.String("token", "", "Slack token")
)

func TestSlack(t *testing.T) {

	if *token == "" {
		fmt.Println("No Token Specified, Skipping")
		t.Skip()
	}

	_, err := NewRTMClient(*token)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		t.Fail()
	}
}
