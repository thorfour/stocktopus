package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thorfour/stocktopus/pkg/auth"
	"github.com/thorfour/stocktopus/pkg/cfg"
	"github.com/thorfour/stocktopus/pkg/stocktopus"
)

const (
	ephemeral = "ephemeral"
	inchannel = "in_channel"
)

var (
	port      = flag.Int("p", 443, "port to serve on")
	notls     = flag.Bool("n", false, "turn off TLS")
	debug     = flag.Bool("d", false, "turn on debugging. Disable TLS")
	certCache = flag.String("c", "/cert", "location to store certs")
)

// response is the json struct for a slack response
type response struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func init() {
	flag.Parse()
}

func main() {
	log.Printf("Starting server on port %v", *port)

	tlsOff := *debug || *notls
	if !tlsOff {
		log.Printf("Serving TLS for host %s", cfg.AllowedHost)
		log.Printf("Storing certs in %s", *certCache)
	}

	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler()) // start prometheus endpoint

	run(*port, tlsOff, *certCache, r)
}

func run(p int, tlsOff bool, certDir string, router *mux.Router) {

	router.HandleFunc("/v1", handler)
	router.HandleFunc("/auth", auth.Dummy())

	if tlsOff { // no TLS

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", p), router))

	} else {

		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(cfg.AllowedHost),
			Cache:      autocert.DirCache(certDir),
			Email:      cfg.SupportEmail,
		}

		srv := &http.Server{
			Handler:   router,
			Addr:      fmt.Sprintf(":%v", p),
			TLSConfig: m.TLSConfig(),
		}
		go http.ListenAndServe(":80", m.HTTPHandler(nil))
		log.Fatal(srv.ListenAndServeTLS("", ""))
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

	// Switch to an ephemeral message
	if err != nil {
		r.ResponseType = ephemeral
		r.Text = err.Error()
	}

	b, err := json.Marshal(r)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(b)
	return
}
