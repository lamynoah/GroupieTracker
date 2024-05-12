package webfunc

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var ptitBacConns = ConnSet{}

type PtitBacData struct {
	// RoomName     string
	RoomLink     string
	PtitBacConns ConnSet
	Letter       string
	ArrayLetter  []string
	IsStarted    bool
	UsersInputs  map[string][]string
	Timer        int
	MaxRound     int
	CurrentRound int
	Categories   []string
	CurrentTime  int
	IsDone       bool
}

var arrayRoom = map[int]*PtitBacData{}

func (room *PtitBacData) SendToRoom(msg any) {
	for v := range room.PtitBacConns {
		if err := v.WriteJSON(msg); err != nil {
			log.Println(err)
			return
		}
	}
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected successfully")
	uri := r.RequestURI
	reader(ws, uri)
}
