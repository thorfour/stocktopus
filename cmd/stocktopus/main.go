package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/acme/autocert"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thorfour/stocktopus/pkg/auth"
	"github.com/thorfour/stocktopus/pkg/slack"
	"github.com/thorfour/stocktopus/pkg/stock"
)

var (
	port         = flag.Int("p", 443, "port to serve on")
	notls        = flag.Bool("n", false, "turn off TLS")
	debug        = flag.Bool("d", false, "turn on debugging. Disable TLS")
	certCache    = flag.String("c", "/cert", "location to store certs")
	allowedHost  = flag.String("host", "api.stocktopus.io", "ACME allowed FQDN")
	supportEmail = flag.String("email", "support@stocktopus.io", "ACME support email")

	redisPW      string
	redisAddr    string
	clientID     string
	clientSecret string
)

func init() {
	flag.Parse()
	redisPW = os.Getenv("REDISPW")
	redisAddr = os.Getenv("REDISADDR")
	clientID = os.Getenv("CLIENTID")
	clientSecret = os.Getenv("CLIENTSECRET")
}

func main() {
	log.Printf("Starting server on port %v", *port)

	tlsOff := *debug || *notls
	if !tlsOff {
		log.Printf("Serving TLS for host %s", allowedHost)
		log.Printf("Storing certs in %s", *certCache)
	}

	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler()) // start prometheus endpoint

	run(*port, tlsOff, *certCache, r)
}

func run(p int, tlsOff bool, certDir string, router *mux.Router) {

	s := slack.New(redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPW,
	}),
		&stock.IexWrapper{},
	)

	router.HandleFunc("/v1", s.Handler)
	router.HandleFunc("/auth", auth.Dummy())

	if tlsOff { // no TLS

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", p), router))

	} else {

		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(allowedHost),
			Cache:      autocert.DirCache(certDir),
			Email:      supportEmail,
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
