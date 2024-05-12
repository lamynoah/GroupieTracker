package webfunc

import (
	"GT/connect"
	"fmt"
	"html/template"
	"net/http"
)

func Select(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/selectGame.html")
	temp.Execute(w, nil)
}

func Lobby(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/lobby.html")
	temp.Execute(w, arrayRoom)
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
