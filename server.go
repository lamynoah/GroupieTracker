package main

import (
	_ "GT/api"
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
	http.HandleFunc("/blindtest", webfunc.BlindtestPage)
	http.HandleFunc("/deaftest", webfunc.DeafTestPage)
	http.HandleFunc("/ptitbac", webfunc.PtitbacPage)
	http.HandleFunc("/getTrackID", webfunc.GetTrackID)
	http.HandleFunc("/ws", webfunc.WebSocket)
	http.HandleFunc("/ws/blindtest", webfunc.WebSocket)
	http.HandleFunc("/ws/deafTest", webfunc.WebSocket)
	http.HandleFunc("/ws/ptitBac", webfunc.WebSocket)
	http.HandleFunc("/ws/loading", webfunc.WebSocket)

	//conn, _ := net.Dial("tcp", "google.com:http")
	//fmt.Println(conn.LocalAddr())

	fmt.Println("Listening on port 8080")
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
