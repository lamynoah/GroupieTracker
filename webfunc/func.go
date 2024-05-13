package webfunc

import (
	"GT/api"
	"GT/bdd"
	"GT/connect"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
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
	id                int
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

type BlindTestData struct {
	ID             int
	RoomLink       string
	BlindTestConns syncmap.Map
	IsStarted      bool
	UsersInputs    sync.Map
	Timer          int
	MaxRounds      int
	CurrentRound   int
	CurrentTime    int
	IsDone         bool
}

func (room *PtitBacData) getId() int {
	return room.id
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


func (room *BlindTestData) SendToRoom(msg any) {
	room.BlindTestConns.Range(func(key any, v any) bool {
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

func GetTrackID(w http.ResponseWriter, r *http.Request) {
	trackID := api.RequestSong()
	w.Write([]byte(trackID))
}

// On prends [3].categories  si majorité false == 0 sinon true == 2 et si reponse non unique score == 1
type Validation struct {
	UserName string `json:"username"`
	Input    string `json:"input"`
	Value    bool   `json:"value"`
}

// var test = map[string]bool{`{"username":"noah12","category":"Album","input":"sss"}`: true, `{"username":"noah12","category":"Artiste","input":"ss"}`: false, `{"username": "noah12", "category": "Featuring", "input": "ss"}`: true, `{"username":"noah12","category":"Groupe de musique","input":"ss"}`: true, `{"username":"noah12","category":"Instrument de musique","input":"ss"}`: true}

func AddScoreToPlayer(roomId, userId, number int) {
	_, scores, err := bdd.QueryRoomUsersScores(roomId)
	if err != nil {
		log.Println("AddScoreToPlayer(QueryRoomUsersScores):", err)
	}
	userScore := scores[userId]
	_ = userScore
	db, err := sql.Open("sqlite3", BDDPath)
	if err != nil {
		log.Println("AddScoreToPlayer:", err)
	}
	defer db.Close()
	_, err = db.Exec("UPDATE ROOM_USERS SET score = score + ? WHERE id_user = ? AND id_room = ?", number, userId, roomId)
	if err != nil {
		log.Println("AddScoreToPlayer(Exec):", err)
	}
}

// Recup la cate  puis
func (room *PtitBacData) IsRight(key string) bool {
	nbplayers := lenOfMap(&room.PtitBacConns)
	count := 0
	for _, v := range room.UsersPointsInputs {
		if v[key] {
			count++
		}
		if count >= nbplayers/2 {
			return true
		}
	}
	return false
}

func (room *PtitBacData) CalcPoints() {
	type Inputed struct {
		Username string `json:"username"`
		Category string `json:"category"`
		Input    string `json:"input"`
	}
	type UserCat struct {
		UserId   int    `json:"username"`
		Category string `json:"category"`
	}
	rightKeys := []string{}
	for i := range room.UsersPointsInputs[0] {
		if room.IsRight(i) {
			rightKeys = append(rightKeys, i)
		}
	}

	inputs := map[string][]UserCat{}

	for _, v := range rightKeys {
		temp := Inputed{}
		err := json.Unmarshal([]byte(v), &temp)
		if err != nil {
			log.Println("CalcPoints(for(Unmarshal)):", err)
		} else {
			fmt.Println("unmarshal result:", temp)
		}

		UserId, err := connect.QueryUserId(temp.Username)
		if err != nil {
			log.Println("CalcPoints(for(QueryUserId)):", err)
		}
		tretedName := CleanStr(strings.ToLower(temp.Input))
		inputs[temp.Input] = append(inputs[tretedName], UserCat{UserId, temp.Category})
	}

	// db, err := sql.Open("sqlite3", BDDPath)
	// if err != nil {
	// 	log.Println("CalcPoints(sql.Open):", err)
	// }

	for _, v := range inputs {
		if len(v) > 1 {
			for _, v2 := range v {
				AddScoreToPlayer(room.getId(), v2.UserId, 1)
			}
		} else {
			AddScoreToPlayer(room.getId(), v[0].UserId, 2)
		}
	}
	// db.Close()
}

func CleanStr(str string) string {
	regex := regexp.MustCompile(`\W+`)
	cleanedStr := regex.ReplaceAllString(str, "")
	return cleanedStr
}
