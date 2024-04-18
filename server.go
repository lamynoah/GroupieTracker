package main

import (
	"GT/webfunc"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.HandleFunc("/selectGame", webfunc.Select)
	http.HandleFunc("/signin", webfunc.Signin)
	http.HandleFunc("/createUser", webfunc.CreateUser)
	http.HandleFunc("/login", webfunc.Login)
	http.HandleFunc("/loginUser", webfunc.Connect)
	http.HandleFunc("/settingBacPage", webfunc.SettingBacPage)
	http.HandleFunc("/blindTest", webfunc.BlindTestPage)
	http.HandleFunc("/deafTest", webfunc.DeafTestPage)
	http.HandleFunc("/ptitbac", webfunc.PtitbacPage)
	http.HandleFunc("/", webfunc.HomePage)
	http.HandleFunc("/ws", webfunc.WebSocket)

	//conn, _ := net.Dial("tcp", "google.com:http")
	//fmt.Println(conn.LocalAddr())

	fmt.Println("Listening on port 8080")
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
