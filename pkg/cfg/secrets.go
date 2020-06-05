package cfg

import "os"

var (
	// CLIENTID is the client ID given from Slack during OAUTH setup
	CLIENTID string
	// CLIENTSECRET is the client secret given from Slack during OAUTH setup
	CLIENTSECRET string
)

func init() {
	CLIENTID = os.Getenv("CLIENTID")
	CLIENTSECRET = os.Getenv("CLIENTSECRET")
}
