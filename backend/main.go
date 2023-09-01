package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ztrue/shutdown"
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
	Type   string `json:"type"`
}

// ** Messages types
// - server.join
// - server.chat
// - server.info
// - server.error

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
var server Server = Server{}

func handleShutdown() {
	message := Message{
		Sender: "server",
		Body:   "The server will shutdown in 10 seconds.",
		Type:   "server.warning",
	}
	broadcast(message)
	fmt.Println("Server will shutdown in 10 seconds...")
	time.Sleep(time.Second * 7)
	fmt.Println("Shutting down in 3...")
	message.Body = "Shutting down in 3..."
	broadcast(message)
	time.Sleep(time.Second)
	fmt.Println("2...")
	message.Body = "2..."
	broadcast(message)
	time.Sleep(time.Second)
	fmt.Println("1...")
	message.Body = "1..."
	broadcast(message)
	time.Sleep(time.Second)
	message.Body = "goodbye."
	broadcast(message)
	fmt.Println("goodbye.")
}
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
func usernameAvailable(username string) bool {
	if username == "server" {
		return false
	}
	for _, user := range server.users {
		if user.name == username {
			return false
		}
	}
	return true
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
func handleServerJoin(message Message, c *websocket.Conn) {
	if usernameAvailable(message.Body) == false {
		json, err := json.Marshal(
			Message{
				Sender: "server",
				Body:   fmt.Sprintf("The name '%s' is unavailable. Please join with another name.", message.Sender),
				Type:   "server.error",
			},
		)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = c.WriteMessage(websocket.TextMessage, json)
		if err != nil {
			fmt.Println(err)
			return
		}
		c.Close()
		return
	}
	newUser := User{message.Body, c}
	addUser(newUser)
	c.SetCloseHandler(func(code int, text string) error {
		fmt.Printf("%s disconnected!\n", message.Body)
		removeUser(newUser)
		broadcast(Message{"server", fmt.Sprintf("%s has disconnected!", message.Body), "server.info"})
		return nil
	})
	fmt.Printf("%s connected!\n", message.Body)
	broadcast(Message{"server", fmt.Sprintf("%s is online!", message.Body), "server.info"})
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
		_, message, err := c.ReadMessage()
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
		parsedMess := Message{
			Sender: fmt.Sprint(transmission["sender"]),
			Body:   fmt.Sprint(transmission["body"]),
			Type:   fmt.Sprint(transmission["type"]),
		}
		switch parsedMess.Type {
		case "server.join":
			handleServerJoin(parsedMess, c)
			break
		case "server.chat":
			broadcast(parsedMess)
			break
		default:
			fmt.Printf("UNKNOWN MESSAGE TYPE RECEIVED: '%s'\n", parsedMess.Type)
		}
	}
}
func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hey!")
	})
	http.HandleFunc("/ws", receive)
}
func printUsage() {
	fmt.Println("go run prog <port number>")
}
func main() {
	port := 5005
	if len(os.Args) > 1 {
		num, err := strconv.Atoi(os.Args[1])
		if err != nil {
			printUsage()
			os.Exit(1)
		}
		port = num
	}
	fmt.Printf("Server has started on port %d!\n", port)
	setupRoutes()
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	shutdown.Add(handleShutdown)
	shutdown.Listen()
}
