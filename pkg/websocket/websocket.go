package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Upgrade(write http.ResponseWriter, request *http.Request) (*websocket.Conn, error) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	connection, err := upgrader.Upgrade(write, request, nil)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return connection, nil
}
