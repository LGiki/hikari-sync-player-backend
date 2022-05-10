package ws

import (
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Hub struct {
	RoomPlayState map[string]*PlayState
	Clients       map[string]map[*Client]bool
	Broadcast     chan WebsocketBroadcastMessage
	Register      chan *Client
	Unregister    chan *Client
}

func NewHub() *Hub {
	return &Hub{
		RoomPlayState: make(map[string]*PlayState),
		Broadcast:     make(chan WebsocketBroadcastMessage),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Clients:       make(map[string]map[*Client]bool),
	}
}

func (h *Hub) CountClients() int {
	return len(h.Clients)
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			clients := h.Clients[client.RoomId]
			if clients == nil {
				clients = make(map[*Client]bool)
				h.Clients[client.RoomId] = clients
			}
			h.Clients[client.RoomId][client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client.RoomId]; ok {
				delete(h.Clients[client.RoomId], client)
				close(client.Send)
				if len(h.Clients[client.RoomId]) == 0 {
					delete(h.Clients, client.RoomId)
					delete(h.RoomPlayState, client.RoomId)
				}
			}
		case message := <-h.Broadcast:
			for client := range h.Clients[message.RoomId] {
				select {
				case client.Send <- message.Data:
				default:
					close(client.Send)
					delete(h.Clients[client.RoomId], client)
					if len(h.Clients[client.RoomId]) == 0 {
						delete(h.Clients, client.RoomId)
						delete(h.RoomPlayState, client.RoomId)
					}
				}
			}
		}
	}
}
