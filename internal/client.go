package server

import "github.com/gofiber/websocket/v2"

type Client struct {
	conn *websocket.Conn
	name string
	id   int
}

func NewClient(connection *websocket.Conn, id int) *Client {
	return &Client{
		conn: connection,
		name: GetRandomName(),
		id:   id,
	}
}
