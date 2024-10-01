package ws

import (
	"fmt"
	"time"
)

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

type Hub struct {
	Rooms      map[string]*Room `json:"rooms"`
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

const EMPTY_ROOM_TIMEOUT = 10 * time.Second

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			if _, ok := h.Rooms[cl.RoomID]; ok {
				r := h.Rooms[cl.RoomID]

				if _, ok := r.Clients[cl.ID]; !ok {
					r.Clients[cl.ID] = cl
				}
				h.scheduleDeleteRoom(h.Rooms[cl.RoomID])
			}
		case cl := <-h.Unregister:
			if _, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := h.Rooms[cl.RoomID].Clients[cl.ID]; ok {
					if len(h.Rooms[cl.RoomID].Clients) != 0 {
						h.Broadcast <- &Message{
							Body:     fmt.Sprintf("%s has left the room", cl.Username),
							RoomID:   cl.RoomID,
							Username: cl.Username,
							UserID:   cl.ID,
						}
					}
					delete(h.Rooms[cl.RoomID].Clients, cl.ID)
					close(cl.Message)
					h.scheduleDeleteRoom(h.Rooms[cl.RoomID])
				}
			}
		case m := <-h.Broadcast:
			if _, ok := h.Rooms[m.RoomID]; ok {
				for _, cl := range h.Rooms[m.RoomID].Clients {
					cl.Message <- m
				}
			}
		}
	}
}

func (h *Hub) scheduleDeleteRoom(room *Room) {
	go func() {
		time.AfterFunc(EMPTY_ROOM_TIMEOUT, func() {
			if len(room.Clients) == 0 || h.Rooms[room.ID] == nil {
				delete(h.Rooms, room.ID)
			}
		})
	}()
}
