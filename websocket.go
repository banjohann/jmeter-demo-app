package main

import (
	"bytes"
	"encoding/json"
	"log"
	"text/template"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type Client struct {
	conn *websocket.Conn
	name string
}

func NewClient(connection *websocket.Conn) *Client {
	return &Client{
		conn: connection,
		name: GetRandomName(),
	}
}

type WebSocketServer struct {
	id        string
	clients   map[*Client]bool
	broadcast chan *Message
}

func NewWebSocket() *WebSocketServer {
	return &WebSocketServer{
		id:        uuid.New().String(),
		clients:   make(map[*Client]bool),
		broadcast: make(chan *Message),
	}
}

func (s *WebSocketServer) HandleWebSocket(ctx *websocket.Conn) {

	wsClient := NewClient(ctx)

	s.clients[wsClient] = true
	defer func() {
		delete(s.clients, wsClient)
		ctx.Close()
	}()

	for {
		_, msg, err := ctx.ReadMessage()
		if err != nil {
			log.Println("Read Error:", err)
			break
		}

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Fatalf("Error Unmarshalling")
		}

		message.ClientName = wsClient.name
		s.broadcast <- &message
	}
}

func (s *WebSocketServer) HandleMessages() {
	for {
		msg := <-s.broadcast

		for client := range s.clients {
			err := client.conn.WriteMessage(websocket.TextMessage, getMessageTemplate(msg))
			if err != nil {
				log.Printf("Write  Error: %v ", err)
				client.conn.Close()
				delete(s.clients, client)
			}
		}
	}
}

func getMessageTemplate(msg *Message) []byte {
	tmpl, err := template.ParseFiles("views/message.html")
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}

	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, msg)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}

	return renderedMessage.Bytes()
}
