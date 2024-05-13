package webfunc

import (
	"GT/games"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/syncmap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const BDDPath = "./bdd/table.db"

// var ptitBacConns = ConnSet{}

type PtitBacData struct {
	RoomLink     		string
	PtitBacConns 		syncmap.Map
	Letter       		string
	ArrayLetter  		[]string
	IsStarted    		bool
	UsersInputs  		syncmap.Map
	Timer        		int
	MaxRound     		int
	CurrentRound 		int
	Categories   		[]string
	CurrentTime  		int
	IsDone       		bool
	UsersPointsInputs  	[]map[string][]games.Validation
}

func (room *PtitBacData) SendToRoom(msg any) {
	room.PtitBacConns.Range(func (key any, v any) bool {
		v.(*sync.Mutex).Lock()
		defer v.(*sync.Mutex).Unlock()
		if err := key.(*websocket.Conn).WriteJSON(msg); err != nil {
			log.Println(err)
			return false
		}
		return true
	})
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
