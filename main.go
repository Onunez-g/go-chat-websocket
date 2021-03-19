package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/onunez-g/go-websocket-tut/pkg/websocket"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(p *websocket.Pool, w http.ResponseWriter, r *http.Request) {

	log.Println(r.Host)

	ws, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}

	client := &websocket.Client{
		Conn: ws,
		Pool: p,
	}
	p.Register <- client
	client.Read()
}

func setupRoutes() *mux.Router {
	pool := websocket.NewPool()
	go pool.Start()
	r := mux.NewRouter()
	r.HandleFunc("/", homePage)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsEndpoint(pool, w, r)
	})
	return r
}

func main() {
	fmt.Println("Welcome to websocket in go")
	r := setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", r))
}
