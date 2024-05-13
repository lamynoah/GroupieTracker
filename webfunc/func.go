package webfunc

import (
	"GT/bdd"
	"database/sql"
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
	RoomLink          string
	PtitBacConns      syncmap.Map
	Letter            string
	ArrayLetter       []string
	IsStarted         bool
	UsersInputs       syncmap.Map
	Timer             int
	MaxRound          int
	CurrentRound      int
	Categories        []string
	CurrentTime       int
	IsDone            bool
	UsersPointsInputs [](map[string]bool)
}

func (room *PtitBacData) SendToRoom(msg any) {
	room.PtitBacConns.Range(func(key any, v any) bool {
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

// On prends [3].categories  si majorité false == 0 sinon true == 2 et si reponse non unique score == 1
type Validation struct {
	UserName string `json:"username"`
	Input    string `json:"input"`
	Value    bool   `json:"value"`
}

var test = map[string]bool{`{"username":"noah12","category":"Album","input":"sss"}`: true, `{"username":"noah12","category":"Artiste","input":"ss"}`: false, `{"username": "noah12", "category": "Featuring", "input": "ss"}`: true, `{"username":"noah12","category":"Groupe de musique","input":"ss"}`: true, `{"username":"noah12","category":"Instrument de musique","input":"ss"}`: true}

func AddScoreToPlayer(roomId, userId, number int) {
	_, scores, err := bdd.QueryRoomUsersScores(roomId)
	if err != nil {
		log.Println("AddScoreToPlayer(QueryRoomUsersScores):", err)
	}
	userScore := scores[userId]
	_ = userScore
	db, err := sql.Open("sqlite3", BDDPath)
	if err != nil {
		log.Println("AddScoreToPlayer(Open):", err)
	}
	_, err = db.Exec("UPDATE ROOM_USERS SET score = score + ? WHERE id_user = ? AND id_room = ?", number, userId, roomId)
	if err != nil {
		log.Println("AddScoreToPlayer(Exec):", err)
	}
}


// Recup la cate  puis 
func (room *PtitBacData) isRight(key string) bool {
	arr := []bool{}
	for _, v := range room.UsersPointsInputs{
		arr = append(arr, v[key])
	}
	return 
}