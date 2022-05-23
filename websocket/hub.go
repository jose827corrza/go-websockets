package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, //Que todos puedan acceder a esto
}

type Hub struct {
	Clients    []*Client
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
	}
}

func (hub *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not connect to the websocket", http.StatusBadRequest)
	}
	client := NewClient(hub, socket)
	hub.register <- client
	go client.Write()
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.onConnect(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

func (hub *Hub) onConnect(client *Client) {
	log.Println("Client connected", client.Socket.RemoteAddr())

	hub.mutex.Lock()
	defer hub.mutex.Unlock()
	client.Id = client.Socket.RemoteAddr().String()
	hub.Clients = append(hub.Clients, client)
}

func (hub *Hub) onDisconnect(client *Client) {
	log.Println("Client disconnected", client.Socket.RemoteAddr())

	client.Socket.Close()
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	i := -1

	for j, c := range hub.Clients {
		if c.Id == client.Id {
			i = j
		}
	}
	copy(hub.Clients[i:], hub.Clients[i+1:])
	hub.Clients[len(hub.Clients)-1] = nil
	hub.Clients = hub.Clients[:len(hub.Clients)-1]
}

func (hub *Hub) Broadcast(message interface{}, ignore *Client) {
	data, _ := json.Marshal(message)
	for _, client := range hub.Clients {
		if client != ignore {
			client.Outbound <- data
		}
	}
}
