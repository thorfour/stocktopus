package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/thorfour/stocktopus/pkg/stocktopus"
)

const (
	ephemeral = "ephemeral"
	inchannel = "in_channel"
)

var (
	port  = flag.Int("p", 8088, "port to serve on")
	debug = flag.Bool("d", false, "turn TLS off")
)

// response is the json struct for a slack response
type response struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func main() {
	flag.Parse()
	log.Printf("Starting server on port %v", *port)
	run(*port, *debug)
}

func run(p int, d bool) {
	http.HandleFunc("/v1", handler)

	if d {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", p), nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%v", p), "cert.pem", "key.pem", nil))
	}
}

func handler(resp http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	// Errors are to be send to the user as an ephemeral message
	msg, err := stocktopus.Process(req.Form)
	newReponse(resp, msg, err)
}

func newReponse(resp http.ResponseWriter, message string, err error) {
	r := &response{
		ResponseType: inchannel,
		Text:         message,
	}

	// Swithc to an ephemeral message
	if err != nil {
		r.ResponseType = ephemeral
		r.Text = err.Error()
	}

	b, err := json.Marshal(r)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Write(b)
	return
}
