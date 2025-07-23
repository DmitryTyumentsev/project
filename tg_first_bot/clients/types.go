package clients

import "net/http"

func NewClient(host string) *Client {
	return &Client{
		HttpClient: http.Client{},
		BasePath:   basePath(),
		Host:       host,
	}
}

type Client struct {
	HttpClient http.Client
	BasePath   string
	Host       string
}

type Response struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateID int `json:"update_id"`
	Message  `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	Text      string `json:"text"`
	User      `json:"from"`
	Chat      `json:"chat"`
}

type User struct {
	UserID int `json:"id"`
}

type Chat struct {
	ChatID   int    `json:"id"`
	Username string `json:"username"`
}
