package slack

// PostRequest is json representation of slack post request
type PostRequest struct {
	token       string `json:"token"`
	teamId      string `json:"team_id"`
	teamDomain  string `json:"team_domain"`
	channelId   string `json:"channel_id"`
	channelName string `json:"channel_name"`
	userId      string `json:"user_id"`
	userName    string `json:"user_name"`
	command     string `json:"command"`
	text        string `json:"text"`
	responseUrl string `json:"response_url"`
}
