package slack

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/thorfour/stocktopus/pkg/stock"
	"github.com/thorfour/stocktopus/pkg/stocktopus"
)

const (
	ephemeral = "ephemeral"
	inchannel = "in_channel"
)

// response is the json struct for a slack response
type response struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

// SlashServer is a slack server that handles slash commands
type SlashServer struct {
	s *stocktopus.Stocktopus
}

// New returns a new slash server
func New(kvstore *redis.Client, stocks stock.Lookup) *SlashServer {
	return &SlashServer{
		s: &stocktopus.Stocktopus{
			KVStore:        kvstore,
			StockInterface: stocks,
		},
	}
}

// Handler is a http handler func for processing slack slash requests for stocktopus
func (s *SlashServer) Handler(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	msg, err := s.Process(ctx, req.Form)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	s.newResponse(resp, msg, nil)
}

// Process a slack request
func (s *SlashServer) Process(ctx context.Context, args url.Values) (string, error) {
	text, ok := args["text"]
	if !ok {
		return "", errors.New("Bad request")
	}

	text = strings.Split(strings.ToUpper(text[0]), " ")
	return s.s.Command(ctx, text[0], text[1:], args)
}

// TODO determine ephermeralness of response
func (s *SlashServer) newResponse(resp http.ResponseWriter, message string, err error) {
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
