package webfunc

import (
	"GT/bdd"
	"database/sql"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"golang.org/x/sync/syncmap"
)

func BlindTestPage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/blindtest.html", "./template/websocket.html")
	temp.Execute(w, nil)
}

type Music struct {
	Title  string
	Lyrics string
}

var musics = []Music{{"a", "aaa"}, {"b", "bbb"}, {"c", "ccc"}}

func (m Music) getLyrics() string { return m.Lyrics }

func (m Music) Check(formValue string) bool {
	if formValue == m.Title {
		return true
	}
	return false
}

func DeafTestPage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/deaftest.html", "./template/websocket.html")
	randomMusic := musics[(rand.Int() % len(musics))]
	temp.Execute(w, randomMusic)
}

// MARK: PtitbacPage
func PtitbacPage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/ptitBac.html", "./template/websocket.html")
	r.ParseForm()
	roomId, err := strconv.Atoi(r.FormValue("room"))
	if err != nil {
		log.Println(err)
	}

	roomRow, err := bdd.QueryRoom(roomId)
	if err != nil {
		log.Println(err)
	}
	userId := getUserIdFromPage(r)

	room := arrayRoom[roomId]
	temp.Execute(w, struct {
		Letter       string
		IsCreator    bool
		RoomId       int
		Categories   []string
		Time         int
		CurrentRound int
		MaxRound     int
		IsDone       bool
	}{
		Letter:       room.Letter,
		IsCreator:    roomRow.Created_by == userId,
		RoomId:       roomId,
		Categories:   room.Categories,
		Time:         room.CurrentTime,
		CurrentRound: room.CurrentRound,
		MaxRound:     room.MaxRound,
		IsDone:       room.IsDone,
	})
}

func SettingBacPage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/settingPagesPtitBac.html")
	temp.Execute(w, nil)
}

func SettingDeaftest(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/settingDeaftest.html")
	temp.Execute(w, nil)
}

func Loading(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/loading.html", "./template/websocket.html")
	r.ParseForm()
	roomId, err := strconv.Atoi(r.FormValue("room"))
	userId := getUserIdFromPage(r)
	db, err := sql.Open("sqlite3", BDDPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	query := "SELECT max_player FROM ROOMS WHERE id = ?"
	row := db.QueryRow(query, roomId)

	var maxPlayer int
	err = row.Scan(&maxPlayer)
	if err != nil {
		http.Redirect(w, r, "/lobby", http.StatusNotFound)
		return
	}

	if lenOfMap(&arrayRoom[roomId].PtitBacConns) > maxPlayer || arrayRoom[roomId].IsStarted {
		http.Redirect(w, r, "/lobby", http.StatusFound)
		return
	}

	room, err := bdd.QueryRoom(roomId)
	if err != nil {
		log.Println(err)
	}

	bdd.InsertRoomsUser(roomId, userId, 0)

	temp.Execute(w, room.Created_by == userId)
}

func lenOfMap(v *syncmap.Map) int {
	i := 0
	v.Range(func(key any, v any) bool {
		i++
		return true
	})
	return i
}

func LoadingPageDeafTest(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/loadingPageDeaftest.html", "./template/websocket.html")
	r.ParseForm()
	roomId, err := strconv.Atoi(r.FormValue("room"))
	userId := getUserIdFromPage(r)
	db, err := sql.Open("sqlite3", BDDPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	query := "SELECT max_player FROM ROOMS WHERE id = ?"
	row := db.QueryRow(query, roomId)

	var maxPlayer int
	err = row.Scan(&maxPlayer)
	if err != nil {
		http.Redirect(w, r, "/lobbyDeaftest", http.StatusNotFound)
		return
	}

	if lenOfMap(&arrayRoomDeaftest[roomId].DeafTestConns) > maxPlayer || arrayRoomDeaftest[roomId].IsStarted {
		http.Redirect(w, r, "/lobbyDeaftest", http.StatusFound)
		return
	}

	room, err := bdd.QueryRoom(roomId)
	if err != nil {
		log.Println(err)
	}

	bdd.InsertRoomsUser(roomId, userId, 0)

	temp.Execute(w, room.Created_by == userId)
}
