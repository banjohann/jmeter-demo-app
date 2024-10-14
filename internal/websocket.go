package server

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type WebSocketServer struct {
	id        string
	clients   map[int]*Client
	broadcast chan *Message
	mu        sync.Mutex
}

func NewWebSocket() *WebSocketServer {
	return &WebSocketServer{
		id:        uuid.New().String(),
		clients:   make(map[int]*Client),
		broadcast: make(chan *Message),
	}
}

func (s *WebSocketServer) getLastClientId() int {
	var lastFound = 0

	for clientId := range s.clients {
		if clientId >= lastFound {
			lastFound = clientId + 1
		}
	}

	return lastFound
}

func (s *WebSocketServer) HandleWebSocket(ctx *websocket.Conn) {
	clientId := s.getLastClientId()
	wsClient := NewClient(ctx, clientId)

	s.clients[clientId] = wsClient

	msgConnected := &Message{
		ClientName: wsClient.name,
		Text:       "Connected Successfully",
		Type:       1,
	}

	err := writeMessage(msgConnected, wsClient.conn)
	if err != nil {
		log.Println(err.Error())
	}

	defer func() {
		delete(s.clients, clientId)
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
			log.Println("Error Unmarshalling message")
			break
		}

		message.ClientName = wsClient.name
		log.Println(message)
		s.broadcast <- &message
	}
}

func writeMessage(msg *Message, connection *websocket.Conn) error {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = connection.WriteMessage(websocket.TextMessage, jsonMsg)
	if err != nil {
		return err
	}

	return nil
}

func (s *WebSocketServer) PublishMessageAllClients(msg *Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshalling message")
		return
	}

	for clientId := range s.clients {
		client := s.clients[clientId]

		err = client.conn.WriteMessage(websocket.TextMessage, jsonMsg)
		if err != nil {
			log.Printf("Write  Error: %v ", err)
			client.conn.Close()
			delete(s.clients, clientId)
		}
	}
}

func (s *WebSocketServer) HandleMessages() {
	for {
		msg := <-s.broadcast
		s.PublishMessageAllClients(msg)
	}
}
