package main

import (
	"fmt"
	"net/http"

	"github.com/dev-el-op/go-websocket/pkg/websocket"
)

func serveWS(pool *websocket.Pool, writer http.ResponseWriter, request *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")

	connection, err := websocket.Upgrade(writer, request)

	if err != nil {
		fmt.Fprintf(writer, "%+v\n", err)
	}

	client := &websocket.Client{
		Connection: connection,
		Pool:       pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/websocket", func(writer http.ResponseWriter, request *http.Request) {
		serveWS(pool, writer, request)
	})
}

func main() {
	setupRoutes()
	http.ListenAndServe(":9000", nil)
	fmt.Println("The websocket server is running on port 9000")
}
