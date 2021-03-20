package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	Type int  `json:"type"`
	Body Body `json:"body"`
}

type Body struct {
	From string `json:"from"`
	To   string `json:"to"`
	Msg  string `json:"msg"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()
	for {
		var body Body
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := json.Unmarshal(p, &body); err != nil {
			fmt.Print(string(p))
			fmt.Errorf("Error: %s", err.Error())
			continue
		}
		message := Message{Type: messageType, Body: body}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
