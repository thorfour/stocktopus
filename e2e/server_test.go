//+build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"testing"
)

// Response is the json struct for a slack response
type Response struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

type Command struct {
	cmd      string
	response string
	match    string
}

func sendCommand(t *testing.T, c Command) string {
	form := url.Values{
		"text":    {c.cmd},
		"user_id": {"test"},
		"token":   {"1234abc"},
	}

	body := bytes.NewBufferString(form.Encode())
	resp, err := http.Post("http://stocktopus:8080/v1", "application/x-www-form-urlencoded", body)
	if err != nil {
		t.Fatalf("post request failed: %v", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read resp body: %v", err)
	}

	r := new(Response)
	if err := json.Unmarshal(b, r); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if r.ResponseType != c.response {
		t.Errorf("unexpected response type (%s): %s", r.ResponseType, r.Text)
	}

	if c.match != "" {
		ok, err := regexp.MatchString(c.match, r.Text)
		if !ok || err != nil {
			t.Errorf("response does not match expectation: %v", r.Text)
		}
	}

	return r.Text
}

func TestSimpleCommands(t *testing.T) {
	commands := []Command{
		{"amd goog", "in_channel", ""},
		{"news amd", "in_channel", ""},
		{"deposit 100000", "ephemeral", "New Balance: 100000"},
		{"buy amd 1", "ephemeral", "Done"},
		{"sell amd 1", "ephemeral", "Done"},
		{"reset", "ephemeral", "New Balance: 0"},
	}
	for _, c := range commands {
		t.Run(c.cmd, func(t *testing.T) {
			sendCommand(t, c)
		})
	}
}
