// Package chatserver implements a chat server with channels for Challenge 8.
package chatserver

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

// Client represents a connected chat client
type Client struct {
	username     string
	messages     chan string
	mu           sync.Mutex
	disconnected bool
}

// Send sends a message to the client (non-blocking, thread-safe).
// Empty messages are ignored so that Receive can use "" to signal connection closed.
func (c *Client) Send(message string) {
	if message == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.disconnected {
		return
	}
	select {
	case c.messages <- message:
	default:
		log.Printf("chat server: dropped message for client %s (channel full)", c.username)
	}
}

// Receive returns the next message for the client (blocking).
// Returns "" when the connection is closed; empty messages are never sent.
func (c *Client) Receive() string {
	msg, ok := <-c.messages
	if !ok {
		return ""
	}
	return msg
}

// ChatServer manages client connections and message routing
type ChatServer struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer {
	return &ChatServer{clients: make(map[string]*Client)}
}

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error) {
	if username == "" {
		return nil, ErrUsernameEmpty
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.clients[username]; exists {
		return nil, ErrUsernameAlreadyTaken
	}
	client := &Client{
		username: username,
		messages: make(chan string, 256),
	}
	s.clients[username] = client
	return client, nil
}

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client) {
	s.mu.Lock()
	delete(s.clients, client.username)
	s.mu.Unlock()

	client.mu.Lock()
	defer client.mu.Unlock()
	if client.disconnected {
		return
	}
	client.disconnected = true
	close(client.messages)
}

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string) {
	if sender == nil || message == "" {
		return
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	formatted := fmt.Sprintf("[%s]: %s", sender.username, message)
	for _, client := range s.clients {
		client.Send(formatted)
	}
}

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	if message == "" {
		return ErrMessageEmpty
	}
	sender.mu.Lock()
	disconnected := sender.disconnected
	sender.mu.Unlock()
	if disconnected {
		return ErrClientDisconnected
	}

	s.mu.RLock()
	recipientClient, exists := s.clients[recipient]
	s.mu.RUnlock()
	if !exists {
		return ErrRecipientNotFound
	}

	recipientClient.Send(fmt.Sprintf("[%s -> %s]: %s", sender.username, recipient, message))
	return nil
}

// Common errors that can be returned by the Chat Server
var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrUsernameEmpty        = errors.New("username cannot be empty")
	ErrMessageEmpty         = errors.New("message cannot be empty")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
)
