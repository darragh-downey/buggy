package websocket

import "fmt"

type Pool struct {
	// Register channel will send out `New User Joined...` to all clients within this pool
	Register chan *Client
	// Unregister will unregister a user and notify the pool when they disconnect
	Unregister chan *Client
	// Clients is a map of clients to a boolean which indicates active/inactive status
	Clients map[*Client]bool
	// Broadcast a channel which when passed a message, will loop through all clients in the pool and send the message through the socket connection
	Broadcast chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of connection pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
			}
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of connection pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User disconnected..."})
			}
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in pool")
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Printf("E: %v", err)
					return
				}
			}
		}
	}
}
