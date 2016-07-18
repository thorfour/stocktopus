package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/websocket"
)

const (
	slackApiUrl = string("https://api.slack.com/")
	rtmstart    = string("https://slack.com/api/rtm.start")
	messageType = string("message")
)

// Standard slack message format
type messageRx struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	User    string `json:"user"`
	Text    string `json:"text"`
	TS      string `json:"ts"`
	// TODO does not include edits or subtypes
}

// Standard slack message format
type messageTx struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type SlackRTMClient struct {
	ws           *websocket.Conn
	id           string
	msg          messageRx
	send         messageTx
	sendSequence uint64
}

func NewRTMClient(token string) (*SlackRTMClient, error) {

	// Request an rtm session
	socketUrl, id, err := rtmStart(token)
	if err != nil {
		return nil, err
	}

	// Connect to the rtm socket
	ws, err := websocket.Dial(socketUrl, "", slackApiUrl)
	if err != nil {
		return nil, err
	}

	return &SlackRTMClient{ws, id, messageRx{}, messageTx{}, 0}, nil
}

func rtmStart(tok string) (socket string, id string, err error) {

	url := fmt.Sprintf("%v?token=%v&no_unreads=true", rtmstart, tok)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("Request Failed: %v", resp.StatusCode)
	}

	defer resp.Body.Close()
	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// Json format
	type self struct {
		Id string `json:"id"`
	}
	type rtmStartResp struct {
		Ok    bool   `json:"ok"`
		Error string `json:"error"`
		Url   string `json:"url"`
		Self  self   `json:"self"`
	}

	// Parse the response
	var data rtmStartResp
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return "", "", err
	}

	if !data.Ok {
		return "", "", fmt.Errorf(data.Error)
	}

	return data.Url, data.Self.Id, nil
}

// Keep reciving data until a message for the requested id is obtained
func (s *SlackRTMClient) Receive() (string, error) {

	checkMsg := func() error {
		if err := websocket.JSON.Receive(s.ws, &s.msg); err != nil {
			return err
		}
		return nil
	}

	var err error
	for err = checkMsg(); err == nil && s.msg.User != s.id && s.msg.Type != messageType; {
		err = checkMsg()
	}

	text := strings.TrimLeft(s.msg.Text, fmt.Sprintf("<@%v>: ", s.id))
	return text, nil
}

// Sends a response to the same channel as previous received message
func (s *SlackRTMClient) Send(m string) error {

	// Setup the send message
	s.send.Channel = s.msg.Channel
	s.sendSequence++
	s.send.Id = s.sendSequence
	s.send.Text = m
	s.send.Type = messageType

	return websocket.JSON.Send(s.ws, &s.send)
}
