package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
var connections []*websocket.Conn = []*websocket.Conn{}

func receive(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	connections = append(connections, c)
	fmt.Println("client connected")
	defer c.Close()
	for {
		mc, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println(c.RemoteAddr(), c.LocalAddr())
			fmt.Println("27:", err)
			return
		}
		fmt.Printf("client says: %s\n", message)
		for _, client := range connections {
			err = client.WriteMessage(mc, message)
			if err != nil {
				fmt.Println("Failed to broadcast message")
				fmt.Println(err)
			}
		}
		// err = c.WriteMessage(mc, message) ( for echo testing )
		if err != nil {
			fmt.Println("write:", err)
			break
		}
	}
}

func setup_routes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hey!")
	})
	http.HandleFunc("/ws", receive)
}
func main() {
	fmt.Println("Server has started on port 5005!")
	setup_routes()
	http.ListenAndServe(":5005", nil)
}
