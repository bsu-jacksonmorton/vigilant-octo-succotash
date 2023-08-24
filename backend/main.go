package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	users []User
	lock  sync.Mutex
}
type User struct {
	name string
	conn *websocket.Conn
}
type Group struct {
	name  string
	users []User
}
type Message struct {
	Sender string `json:"sender"`
	Body   string `json:"body"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
var server Server = Server{}

func removeUser(user User) {
	server.lock.Lock()
	defer server.lock.Unlock()
	for i, currUser := range server.users {
		if currUser.name == user.name {
			server.users = append(server.users[:i], server.users[i+1:]...)
		}
	}
}
func addUser(user User) {
	server.lock.Lock()
	defer server.lock.Unlock()
	server.users = append(server.users, user)
}
func usernameTaken(username string) bool {
	for _, user := range server.users {
		if user.name == username {
			return true
		}
	}
	return false
}
func broadcast(message Message) {
	json, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, user := range server.users {
		err = user.conn.WriteMessage(1, json)
		if err != nil {
			fmt.Println("Failed to broadcast message")
			fmt.Println(err)
		}
	}
}
func receive(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("client connected")
	defer c.Close()
	for {
		mc, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		// Parse message
		fmt.Printf("client says: %s\n", message)
		var transmission map[string]interface{}
		err = json.Unmarshal([]byte(message), &transmission)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// HANDLE "server.join"
		if transmission["type"] == "server.join" {
			// Create User struct for that client and add them to the connections
			// check if username taken
			username := fmt.Sprint(transmission["username"])
			if usernameTaken(username) {
				/*
					TODO: Add better handling for username
				*/
				c.Close()
				return
			}
			newUser := User{username, c}
			addUser(newUser)
			c.SetCloseHandler(func(code int, text string) error {
				fmt.Printf("%s disconnected!\n", newUser.name)
				removeUser(newUser)
				broadcast(Message{"server", fmt.Sprintf("%s has disconnected!", newUser.name)})
				return nil
			})
			fmt.Printf("%s connected!\n", newUser.name)
			res := Message{"server", fmt.Sprintf("%s is online!", transmission["username"])}
			json, err := json.Marshal(res)
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, user := range server.users {
				fmt.Println(mc)
				err = user.conn.WriteMessage(mc, json)
				if err != nil {
					fmt.Println("Failed to broadcast message")
					fmt.Println(err)
				}
			}
		} else {
			for _, user := range server.users {
				err = user.conn.WriteMessage(mc, message)
				if err != nil {
					fmt.Println("Failed to broadcast message")
					fmt.Println(err)
				}
			}
		}
	}
}
func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hey!")
	})
	http.HandleFunc("/ws", receive)
}
func main() {
	fmt.Println("Server has started on port 5005!")
	setupRoutes()
	http.ListenAndServe(":5005", nil)
}
