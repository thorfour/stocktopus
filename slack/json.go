package slack

type PostRequest struct {
	token        string `json:"token"`
	team_id      string `json:"team_id"`
	team_domain  string `json:"team_domain"`
	channel_id   string `json:"channel_id"`
	channel_name string `json:"channel_name"`
	user_id      string `json:"user_id"`
	user_name    string `json:"user_name"`
	command      string `json:"command"`
	text         string `json:"text"`
	response_url string `json:"response_url"`
}
