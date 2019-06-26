//+build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

// Response is the json struct for a slack response
type Response struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func TestSimpleQuote(t *testing.T) {
	form := url.Values{
		"text": {"amd"},
	}

	time.Sleep(time.Second)

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

	if r.ResponseType != "in_channel" {
		t.Errorf("unexpected response type: %v", r.ResponseType)
	}
}
