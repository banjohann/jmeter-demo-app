package main

import (
	"bytes"
	"encoding/json"
	"log"
	"sync"
	"text/template"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
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
	mu        sync.Mutex
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
		wsClient.conn.WriteMessage(websocket.TextMessage, getNewForm())
		s.broadcast <- &message
	}
}

func (s *WebSocketServer) PublishMessage(msg *Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for client := range s.clients {
		err := client.conn.WriteMessage(websocket.TextMessage, getMessageTemplate(msg))
		if err != nil {
			log.Printf("Write  Error: %v ", err)
			client.conn.Close()
			delete(s.clients, client)
		}
	}
}

func (s *WebSocketServer) HandleMessages() {
	for {
		msg := <-s.broadcast

		s.PublishMessage(msg)
	}
}

func getNewForm() []byte {
	tmpl, err := template.ParseFiles("views/form.html")
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}

	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, "")
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}

	return renderedMessage.Bytes()
}

func getMessageTemplate(msg *Message) []byte {
	p := bluemonday.NewPolicy()

	message := p.Sanitize(msg.Text)
	msg.Text = message

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
