package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-management/service/jwt"

	"github.com/gorilla/mux"
)

type Broker struct {
	// users is a map where the key is the user ID and the value
	users map[string]chan []byte
	// actions is a channel of functions to cal in the brokers
	// go routine. The broker executes everything in that single
	// go routine to avoid data races.
	actions  chan func()
	jwtMaker *jwt.JWTMaker
}

type NotifMessage struct {
	Name     string `json:"name"`
	Message  string `json:"message"`
	RoomID   string `json:"roomID"`
	RoomName string `json:"roomName"`
}

func NewBroker(secretKey string) *Broker {
	return &Broker{
		users:    make(map[string]chan []byte),
		actions:  make(chan func()),
		jwtMaker: jwt.NewJWTMaker(secretKey),
	}
}

func (b *Broker) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/notifications/stream/{id}", b.stream)
	router.HandleFunc("/notifications/send/{id}", b.broadcastMessage).Methods("POST")
	go b.listen()
}

func (b *Broker) listen() {
	for action := range b.actions {
		action()
	}
}

func (b *Broker) addUserChan(id string, ch chan []byte) {
	b.actions <- func() {
		b.users[id] = ch
	}
}

// removeUserChan removes a channel for a user with the given ID
func (b Broker) removeUserChan(id string, ch chan []byte) {
	go func() {
		for range ch {
		}
	}()
	b.actions <- func() {
		delete(b.users, id)
		close(ch)
	}
}

func (b *Broker) stream(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	flusher := w.(http.Flusher)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ch := make(chan []byte)
	if _, ok := b.users[id]; !ok {
		b.addUserChan(id, ch)
	}
	defer b.removeUserChan(id, ch)
	for {
		select {
		case <-r.Context().Done():
			return
		case m := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", m)
			flusher.Flush()
		}
	}
}

func (b *Broker) broadcastMessage(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var message NotifMessage
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	j, err := json.Marshal(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b.sendToUser(id, j)
}

func (b *Broker) sendToUser(id string, data []byte) {
	b.actions <- func() {
		b.users[id] <- data
	}
}
