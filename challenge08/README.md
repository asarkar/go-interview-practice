# Challenge 8: Chat Server with Channels

## Problem Statement

Implement a simple chat server using Go channels and goroutines. The chat server should allow multiple clients to connect, broadcast messages to all clients, and support private messaging between clients.

## Requirements

1. Implement a `ChatServer` struct that manages connections and message routing:
   - Add and remove clients
   - Broadcast messages to all clients
   - Route private messages between specific clients
   - Handle disconnections gracefully

2. Implement a `Client` struct that represents a connected client:
   - Unique username
   - Incoming message channel
   - Outgoing message channel
   - Connection status

3. Use channels to manage message flow between clients and the server.

4. Implement concurrency using goroutines for handling multiple clients simultaneously.

5. Create test cases that simulate multiple clients connecting/disconnecting and exchanging messages.

## Function Signatures

```go
// ChatServer manages client connections and message routing
type ChatServer struct {
    // Your implementation here
}

// NewChatServer creates a new chat server instance
func NewChatServer() *ChatServer

// Connect adds a new client to the chat server
func (s *ChatServer) Connect(username string) (*Client, error)

// Disconnect removes a client from the chat server
func (s *ChatServer) Disconnect(client *Client)

// Broadcast sends a message to all connected clients
func (s *ChatServer) Broadcast(sender *Client, message string)

// PrivateMessage sends a message to a specific client
func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error

// Client represents a connected chat client
type Client struct {
    // Your implementation here
}

// Send sends a message to the client
func (c *Client) Send(message string)

// Receive returns the next message for the client (blocking)
func (c *Client) Receive() string
```

## Test Cases

Your implementation should handle the following test scenarios:

1. Multiple clients connecting to the server
2. Broadcasting messages to all clients
3. Sending private messages between specific clients
4. Clients disconnecting and reconnecting
5. Error handling for invalid operations (e.g., sending to non-existent clients)
6. Concurrent operations (multiple clients sending/receiving simultaneously)

## Constraints

- All operations must be thread-safe
- The chat server should handle at least 100 concurrent clients
- Messages should be delivered in the order they are sent
- Handle errors gracefully and provide meaningful error messages
- Test cases should not result in deadlocks or race conditions

## Instructions

- **Fork** the repository.
- **Clone** your fork to your local machine.
- **Create** a directory named after your GitHub username inside `challenge-8/submissions/`.
- **Copy** the `solution-template.go` file into your submission directory.
- **Implement** the required structs and methods.
- **Test** your solution locally by running the test file.
- **Commit** and **push** your code to your fork.
- **Create** a pull request to submit your solution.

## Testing Your Solution Locally

Run the following command in the `challenge-8/` directory:

```bash
go test -v
```

## Advanced Bonus Challenges

For those seeking extra challenges:

1. Implement a message history feature that allows clients to retrieve the last N messages when they connect
2. Add support for "chat rooms" where clients can join/leave specific rooms
3. Implement a timeout feature that disconnects idle clients after a specified period
4. Add support for message types (e.g., text, image references, system notifications) 