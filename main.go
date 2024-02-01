package main

import (
	"fmt"
	"net/http"

	"github.com/dev-el-op/go-websocket/helpers"
	"github.com/dev-el-op/go-websocket/pkg/websocket"
)

func serveWS(pool *websocket.Pool, writer http.ResponseWriter, request *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")

	channel := request.URL.Query().Get("channel")
	authToken := request.URL.Query().Get("token")

	if helpers.IsValidToken(channel, authToken) {

		connection, err := websocket.Upgrade(writer, request)

		if err != nil {
			fmt.Fprintf(writer, "%+v\n", err)
		}

		client := &websocket.Client{
			Connection: connection,
			Pool:       pool,
			Channel:    channel,
			AuthToken:  authToken,
		}

		pool.Register <- client
		client.Read()

	} else {
		http.Error(writer, "Invalid channel or token", http.StatusUnauthorized)
	}
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/websocket", func(writer http.ResponseWriter, request *http.Request) {
		serveWS(pool, writer, request)
	})

	http.HandleFunc("/add-channel", func(writer http.ResponseWriter, request *http.Request) {
		helpers.AddChannel(writer, request)
	})
}

func main() {
	setupRoutes()
	http.ListenAndServe(":9000", nil)
	fmt.Println("The websocket server is running on port 9000")
}
