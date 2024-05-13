package webfunc

import (
	"GT/bdd"
	"GT/connect"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func Select(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/selectGame.html")
	temp.Execute(w, nil)
}

func Lobby(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/lobby.html")
	arrayNamedRoom := map[string]*PtitBacData{}
	for i, v := range arrayRoom {
		row, err := bdd.QueryRoom(i)
		if err != nil {
			log.Println(err)
		}
		arrayNamedRoom[row.Name] = v
	}
	temp.Execute(w, arrayNamedRoom)
}

func LobbyDeaftest(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/lobbyDeaftest.html")
	arrayNamedRoom := map[string]*DeafTestData{}
	for i, v := range arrayRoomDeaftest {
		row, err := bdd.QueryRoom(i)
		if err != nil {
			log.Println(err)
		}
		arrayNamedRoom[row.Name] = v
	}
	temp.Execute(w, arrayNamedRoom)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/home.html", "./template/websocket.html")
	temp.Execute(w, nil)
}

func Signin(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/signin.html", "./template/signin.html")
	temp.Execute(w, dataError)
	dataError.Error = ""
}

var dataError = struct {
	Error string
}{""}

func Login(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/login.html")
	fmt.Println(connect.Islogin(r))
	temp.Execute(w, dataError)
	dataError.Error = ""
}

func Score(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/score.html", "./template/scoreboard.html")

	r.ParseForm()
	roomId := getRoomIdFromPage(r)
	scores, _, err := bdd.QueryRoomUsersScores(roomId)
	if err != nil {
		log.Fatal(err)
	}
	temp.Execute(w, scores)
}
