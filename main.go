package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/onunez-g/go-websocket-tut/pkg/websocket"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func clientsEndpoint(p *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	clients := p.GetClients()
	response, err := json.Marshal(clients)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(response)
}

func wsEndpoint(p *websocket.Pool, w http.ResponseWriter, r *http.Request) {

	log.Println(r.Host)
	ws, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}
	queryString := r.URL.Query().Get("id")
	if queryString == "" {
		queryString = fmt.Sprintf("Anonymous-%d", rand.Intn(1000))
	}
	client := &websocket.Client{
		ID:   queryString,
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
	r.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		clientsEndpoint(pool, w, r)
	})
	return r
}

func main() {
	fmt.Println("Welcome to websocket in go")
	r := setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", r))
}
