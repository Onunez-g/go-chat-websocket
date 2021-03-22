package websocket

import "fmt"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[string]*Client
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan Message),
	}
}

func (p *Pool) Start() {
	for {
		select {
		case client := <-p.Register:
			p.Clients[client.ID] = client
			fmt.Println("Size of Connection Pool: ", len(p.Clients))
			for _, c := range p.Clients {
				fmt.Println(c)
				c.Conn.WriteJSON(Message{Type: 1, Body: Body{
					From: "",
					To:   "all",
					Msg:  fmt.Sprintf("%s joined the chat!", client.ID),
				}})
			}
			break
		case client := <-p.Unregister:
			delete(p.Clients, client.ID)
			fmt.Println("Size of Connection Pool: ", len(p.Clients))
			for _, c := range p.Clients {
				c.Conn.WriteJSON(Message{Type: 1, Body: Body{
					From: "",
					To:   "all",
					Msg:  fmt.Sprintf("%s disconnected", client.ID),
				}})
			}
			break
		case message := <-p.Broadcast:

			if message.Body.To == "all" || message.Body.To == "" {
				fmt.Println("Sending message to all clients in Pool")
				for _, c := range p.Clients {
					if err := c.Conn.WriteJSON(message); err != nil {
						fmt.Println(err)
						return
					}
				}
			} else {
				fmt.Printf("Sending message to %s", message.Body.To)
				if err := p.Clients[message.Body.To].Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
				if err := p.Clients[message.Body.From].Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("sent")
			}
		}
	}
}

func (p *Pool) GetClients() []string {
	keys := make([]string, 0, len(p.Clients))
	for k := range p.Clients {
		keys = append(keys, k)
	}
	return keys
}
