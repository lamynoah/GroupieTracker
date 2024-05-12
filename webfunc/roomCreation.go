package webfunc

import (
	"GT/bdd"
	"GT/games"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// MARK: CreateRoom
func CreateRoom(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	UserId := getUserIdFromPage(r)
	max_player, err := strconv.Atoi(r.FormValue("playersNumber"))
	if err != nil {
		log.Fatal(err)
	}

	timer, err := strconv.Atoi(r.FormValue("timerSeconds"))
	if err != nil {
		log.Fatal(err)
	}

	name := r.FormValue("name")
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	room := games.ROOM{
		Created_by:  UserId,
		Max_players: max_player,
		Name:        name,
		Id_game:     3,
	}
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	bdd.InsertRooms(room.Created_by, room.Max_players, room.Name, room.Id_game)
	roomID, err := bdd.GetRoomID(room.Name)
	if err != nil {
		log.Fatal(err)
	}

	maxRound, err := strconv.Atoi(r.FormValue("roundsNumber"))
	if err != nil {
		log.Fatal(err)
	}

	letters := []string{}
	letter := games.GenerateUniqueLetters(&letters)

	fmt.Println("catJSON :", r.FormValue("JSON"))
	var categories []string
	if err = json.Unmarshal([]byte(r.FormValue("JSON")), &categories); err != nil {
		log.Println(err)
	}
	fmt.Println(categories)

	//MARK: Init Room
	arrayRoom[roomID] = &PtitBacData{
		RoomLink:     "?room=" + fmt.Sprint(roomID),
		PtitBacConns: ConnSet{},
		Letter:       letter,
		ArrayLetter:  letters,
		IsStarted:    false,
		UsersInputs:  make(map[string][]string),
		Timer:        timer,
		MaxRound:     maxRound,
		CurrentRound: 1,
		Categories:   categories,
		CurrentTime:  timer,
	}

	http.Redirect(w, r, "/loadingPage?room="+strconv.Itoa(roomID), http.StatusFound)
}
