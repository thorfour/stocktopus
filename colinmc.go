//+build !AWS

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/colinmc/slack"
	"github.com/colinmc/stock"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(os.Stderr, "usage: colinmc: slack-bot-token\n")
		return
	}

	// Open connection
	slackBot, err := slack.NewRTMClient(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	for {
		msg, err := slackBot.Receive()
		if err != nil {
			log.Fatal(err)
		}

		if len(msg) != 0 {
			quote, err := stock.GetQuoteGoogle(msg)
			if err != nil {
				continue
			}

			// Post the quote
			slackBot.Send(quote)
		}
	}
}
