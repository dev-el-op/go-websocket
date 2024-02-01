package websocket

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "./tokens.db"

func init() {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("Could not open the database connection:", err)
		os.Exit(1)
	}
	defer db.Close()

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS tokens (
            channel TEXT PRIMARY KEY,
            token   TEXT
        );
    `)

	if err != nil {
		fmt.Println("Could not create the table:", err)
		os.Exit(1)
	}
}

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
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
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client := range pool.Clients {
				fmt.Println(client)
				client.Connection.WriteJSON(Message{Type: 1, Body: "New User Joined...", Channel: client.Channel})
			}
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client := range pool.Clients {
				client.Connection.WriteJSON(Message{Type: 1, Body: "User Disconnected...", Channel: client.Channel})
			}
		case message := <-pool.Broadcast:
			fmt.Printf("Sending message to all clients in channel %s\n", message.Channel)
			for client := range pool.Clients {

				if client.Channel == message.Channel {
					if err := client.Connection.WriteJSON(message); err != nil {
						fmt.Println(err)
						return
					}
				}
			}
		}
	}
}
