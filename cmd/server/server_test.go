package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestSimpleQuote(t *testing.T) {

	// Run a debug server
	go run(8088, true, ".")

	form := url.Values{
		"text": {"amd"},
	}

	time.Sleep(time.Second)

	body := bytes.NewBufferString(form.Encode())
	resp, err := http.Post("http://localhost:8088/v1", "application/x-www-form-urlencoded", body)
	if err != nil {
		t.Fatalf("post request failed: %v", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read resp body: %v", err)
	}

	r := new(response)
	if err := json.Unmarshal(b, r); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if r.ResponseType != inchannel {
		t.Errorf("unexpected response type: %v", r.ResponseType)
	}
}
