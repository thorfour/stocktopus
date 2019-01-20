package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thorfour/stocktopus/pkg/cfg"
	"github.com/thorfour/stocktopus/pkg/stocktopus"
)

const (
	ephemeral = "ephemeral"
	inchannel = "in_channel"
)

var (
	port      = flag.Int("p", 443, "port to serve on")
	debug     = flag.Bool("d", false, "turn TLS off")
	certCache = flag.String("c", "/cert", "location to store certs")
)

// response is the json struct for a slack response
type response struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func main() {
	flag.Parse()
	log.Printf("Starting server on port %v", *port)
	if !*debug {
		log.Printf("Serving TLS for host %s", cfg.AllowedHost)
		log.Printf("Storing certs in %s", *certCache)
	}
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler()) // start prometheus endpoint

	run(*port, *debug, *certCache, r)
}

func run(p int, d bool, certDir string, router *mux.Router) {

	if d {
		router.HandleFunc("/v1", handler)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", p), router))
	} else {
		router.HandleFunc("/v1", handler)
		hostPolicy := func(ctx context.Context, host string) error {
			if host == cfg.AllowedHost {
				return nil
			}
			return fmt.Errorf("acme/autocert: only %s allowed", cfg.AllowedHost)
		}
		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
			Cache:      autocert.DirCache(certDir),
			Email:      cfg.SupportEmail,
		}
		srv := &http.Server{
			Handler: router,
			Addr:    fmt.Sprintf(":%v", p),
			TLSConfig: &tls.Config{
				GetCertificate: m.GetCertificate,
			},
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

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(b)
	return
}

func hostPolicy(ctx context.Context, host string) error {
	if host == cfg.AllowedHost {
		return nil
	}

	return fmt.Errorf("acme/autocert: only %s hist is allowed", cfg.AllowedHost)
}
