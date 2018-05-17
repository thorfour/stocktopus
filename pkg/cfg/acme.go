package cfg

import "flag"

var (
	// AllowedHost for the ACME protocol for TLS certs
	AllowedHost string
	// SupportEmail email for ACME provider to contact for TLS problems
	SupportEmail string
)

func init() {
	flag.StringVar(&AllowedHost, "host", "api.stocktopus.io", "ACME allowed FQDN")
	flag.StringVar(&SupportEmail, "email", "support@stocktopus.io", "ACME support email")
}
