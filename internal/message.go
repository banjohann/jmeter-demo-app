package server

type Message struct {
	Text       string `json:"text"`
	ClientName string `json:"client_name"`
	Type       int    `json:"type"`
}
