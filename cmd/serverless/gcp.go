//+build GCP,!AWS

package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// Successful command print to stdout, errors and ephermeral messages print to stderr
func main() {

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
		return
	}

	// Expect args(1) to be a url encoded string
	decodedMap, err := url.ParseQuery(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: url.ParseQuery")
		return
	}

	text := decodedMap["text"]
	text = strings.Split(strings.ToUpper(text[0]), " ")

	cmd, ok := cmds[text[0]]
	if !ok { // If there is no cmd mapped, assume it's a ticker and get quotes
		getQuotes(decodedMap["text"][0], decodedMap)
	} else {
		cmd.funcPtr(text, decodedMap)
	}
}
