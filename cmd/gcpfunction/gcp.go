package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/thorfour/stocktopus/pkg/stocktopus"
)

// Successful command print to stdout, errors and ephermeral messages print to stderr
func main() {

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Error: Invalid number arguments")
		return
	}

	// Expect args(1) to be a json map
	simpleMap := make(map[string]string)
	if err := json.Unmarshal([]byte(os.Args[1]), &simpleMap); err != nil {
		fmt.Fprintln(os.Stderr, "Error: json.Unmarshal", err)
		return
	}
	//`{"token":"J3Y6nj4YDGtIp6IICPD4kzmO","team_id":"T0FA8NMKQ","team_domain":"currentandformerhgst","channel_id":"G1SQ4CB5L","channel_name":"privategroup","user_id":"U0FLWC43B","user_name":"thor","command":"/spbeta","text":"amd","response_url":"https://hooks.slack.com/commands/T0FA8NMKQ/269055428499/L9BNxCckPl8cLYQFRzwYyGXO"}`

	// Convert the simple map to a url.Values
	decodedMap := make(map[string][]string)
	for k, v := range simpleMap {
		decodedMap[k] = []string{v}
	}

	msg, err := stocktopus.Process(decodedMap)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	fmt.Println(msg)
}
