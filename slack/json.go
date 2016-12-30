package slack

// PostRequest is json representation of slack post request
type PostRequest struct {
	Token       string `json:"token"`
	TeamId      string `json:"team_id"`
	TeamDomain  string `json:"team_domain"`
	ChannelId   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	UserId      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Command     string `json:"command"`
	Text        string `json:"text"`
	ResponseUrl string `json:"response_url"`
}
