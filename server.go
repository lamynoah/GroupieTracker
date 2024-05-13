package main

import (
	_ "GT/api"
	"GT/bdd"
	"GT/webfunc"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// bdd.CreateUserTable()
	// bdd.CreateRoomsTable()
	// bdd.CreateRoomUsersTable()
	// bdd.CreateGamesTable()
	http.HandleFunc("/", webfunc.HomePage)
	// user managing routes
	http.HandleFunc("/signin", webfunc.Signin)
	http.HandleFunc("/createUser", webfunc.CreateUser)
	http.HandleFunc("/login", webfunc.Login)
	http.HandleFunc("/loginUser", webfunc.Connect)
	http.HandleFunc("/deaftest", webfunc.DeafTestPage)
	http.HandleFunc("/ptitbac", webfunc.PtitbacPage)
	http.HandleFunc("/getTrackID", webfunc.GetTrackID)
	// games routes
	http.HandleFunc("/selectGame", webfunc.Select)
	http.HandleFunc("/blindtest", webfunc.BlindTestPage)
	http.HandleFunc("/settingBacPage", webfunc.SettingBacPage)
	http.HandleFunc("/lobby", webfunc.Lobby)
	http.HandleFunc("/lobbyBlindtest", webfunc.LobbyBlindtest)
	http.HandleFunc("/loadingPage", webfunc.Loading)
	http.HandleFunc("/loadingPageBlindtest", webfunc.LoadingPageBlindtest)
	http.HandleFunc("/createRoom", webfunc.CreateRoom)
	http.HandleFunc("/settingBlindtest", webfunc.SettingBlindtest)
	http.HandleFunc("/createRoomBlindtest", webfunc.CreateRoomBlindtest)
	http.HandleFunc("/score", webfunc.Score)
	// http.HandleFunc("/result", webfunc.Result)
	// websockets routes
	http.HandleFunc("/ws", webfunc.WebSocket)
	http.HandleFunc("/ws/blindtest", webfunc.WebSocket)
	http.HandleFunc("/ws/deafTest", webfunc.WebSocket)
	http.HandleFunc("/ws/ptitBac", webfunc.WebSocket)
	http.HandleFunc("/ws/loading", webfunc.WebSocket)
	http.HandleFunc("/ws/loadingBlindtest", webfunc.WebSocket)
	http.HandleFunc("/ws/result", webfunc.WebSocket)
	//conn, _ := net.Dial("tcp", "google.com:http")
	//fmt.Println(conn.LocalAddr())

	// Reset Rooms & RoomUsers
	bdd.DeleteRooms()
	bdd.DeleteRoomUsers()

	// Create css directory
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
