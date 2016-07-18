package main

import (
	"fmt"
	"log"
	"os"

	"github.com/colinmc/slack"
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
		fmt.Println(msg)
		//fmt.Println(stock.GetQuote(msg))
	}
}
