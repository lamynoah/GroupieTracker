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

var arrayRoom = map[int]*PtitBacData{}
var arrayRoomDeaftest = map[int]*DeafTestData{}

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
	db, err := sql.Open("sqlite3", BDDPath)
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

	// fmt.Println("catJSON :", r.FormValue("JSON"))
	var categories []string
	if err = json.Unmarshal([]byte(r.FormValue("JSON")), &categories); err != nil {
		log.Println(err)
	}
	// fmt.Println(categories)

	//MARK: Init Room
	arrayRoom[roomID] = &PtitBacData{
		id:           roomID,
		RoomLink:     "?room=" + fmt.Sprint(roomID),
		Letter:       letter,
		ArrayLetter:  letters,
		IsStarted:    false,
		Timer:        timer,
		MaxRound:     maxRound,
		CurrentRound: 1,
		Categories:   categories,
		CurrentTime:  timer,
	}

	http.Redirect(w, r, "/loadingPage?room="+fmt.Sprint(roomID), http.StatusFound)
}

func CreateRoomDeafTest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	UserId := getUserIdFromPage(r)

	max_player, err := strconv.Atoi(r.FormValue("playersNumber"))
	if err != nil {
		log.Fatal(err)
	}

	name := r.FormValue("name")
	vInt := parseIntValuesFromPage(r, []string{"timerSeconds", "roundsNumber"})

	db, err := sql.Open("sqlite3", BDDPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	room := games.ROOM{
		Created_by:  UserId,
		Max_players: max_player,
		Name:        name,
		Id_game:     2,
	}
	bdd.InsertRooms(room.Created_by, room.Max_players, room.Name, room.Id_game)
	roomID, err := bdd.GetRoomID(room.Name)
	if err != nil {
		log.Fatal(err)
	}

	arrayRoomDeaftest[roomID] = &DeafTestData{
		Id:           roomID,
		RoomLink:     "?room=" + fmt.Sprint(roomID),
		IsStarted:    false,
		Timer:        vInt[0],
		MaxRound:     vInt[1],
		CurrentRound: 1,
		CurrentTime:  vInt[0],
	}
	fmt.Println("here :", "/loadingPageDeafTest?room="+fmt.Sprint(roomID))
	http.Redirect(w, r, "/loadingPageDeafTest?room="+fmt.Sprint(roomID), http.StatusFound)
}

func parseValuesFromPage(r *http.Request, names []string) []string {
	ans := []string{}
	for _, v := range names {
		ans = append(ans, r.FormValue(v))
	}
	return ans
}

func parseIntValuesFromPage(r *http.Request, names []string) []int {
	ans := []int{}
	for _, v := range names {
		value, err := strconv.Atoi(r.FormValue(v))
		if err != nil {
			log.Fatal(err)
		}
		ans = append(ans, value)
	}
	return ans
}
