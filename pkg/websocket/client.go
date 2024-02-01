package websocket

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID         string
	Connection *websocket.Conn
	Pool       *Pool
	Mutex      sync.Mutex
	Channel    string
	AuthToken  string
}

type Message struct {
	Type    int    `json:"type"`
	Body    string `json:"body"`
	Channel string `json:"channel"`
}

func (client *Client) Read() {
	defer func() {
		client.Pool.Unregister <- client
		client.Connection.Close()
	}()

	for {

		messageType, p, err := client.Connection.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}

		message := Message{Type: messageType, Body: string(p), Channel: client.Channel}

		client.Pool.Broadcast <- message

		fmt.Printf("Message Received: %+v\n", message)

	}

}
