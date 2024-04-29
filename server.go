package main

import (
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
	// games routes
	http.HandleFunc("/selectGame", webfunc.Select)
	http.HandleFunc("/blindTest", webfunc.BlindTestPage)
	http.HandleFunc("/deafTest", webfunc.DeafTestPage)
	http.HandleFunc("/ptitBac", webfunc.PtitbacPage)
	http.HandleFunc("/settingBacPage", webfunc.SettingBacPage)
	http.HandleFunc("/loadingpage", webfunc.Loading)
	// websockets routes
	http.HandleFunc("/ws", webfunc.WebSocket)
	http.HandleFunc("/ws/blindTest", webfunc.WebSocket)
	http.HandleFunc("/ws/deafTest", webfunc.WebSocket)
	http.HandleFunc("/ws/ptitBac", webfunc.WebSocket)

	//conn, _ := net.Dial("tcp", "google.com:http")
	//fmt.Println(conn.LocalAddr())

	fmt.Println("Listening on port 8080")
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
