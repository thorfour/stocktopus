package auth

import (
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/thorfour/stocktopus/pkg/cfg"
)

const (
	oathURL     = "https://slack.com/api/oauth.access"
	encodedType = "application/x-www-form-urlencoded"
)

// Dummy is a dummy function for slack oauth handling
// It does no actual auth but satisfies the oauth workflow
func Dummy() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Temp auth code from slack
		code := req.URL.Query().Get("code")

		postURL, _ := url.Parse(oathURL)
		params := url.Values{}
		params.Add("client_id", cfg.CLIENTID)
		params.Add("client_secret", cfg.CLIENTSECRET)
		params.Add("code", code)
		postURL.RawQuery = params.Encode()

		resp, err := http.Post(postURL.String(), encodedType, nil)
		if err != nil {
			log.WithField("err", err).Error("Failed http post")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte("https://github.com/thorfour/stocktopus"))
	}
}
