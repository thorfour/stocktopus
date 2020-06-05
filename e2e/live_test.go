//+build live

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestLiveRequest(t *testing.T) {

	form := url.Values{
		"text": {"amd"},
	}

	body := bytes.NewBufferString(form.Encode())
	resp, err := http.Post("https://api.stocktopus.io:443/v1", "application/x-www-form-urlencoded", body)
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
