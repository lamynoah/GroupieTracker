package webfunc

import (
	. "GT/BDD"
	. "GT/Connect"
	"database/sql"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func Select(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/selectGame.html")
	temp.Execute(w, nil)
}

func Signin(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/signin.html", "./template/signin.html")
	temp.Execute(w, nil)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	hasedPassword, err := HashPassword(password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = InsertUser(username, email, hasedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/selectGame", http.StatusFound)
}
func Login(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/login.html")
	temp.Execute(w, nil)
}

func Connect(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	usernameOrEmail := r.FormValue("usernameOrEmail")
	password := r.FormValue("password")
	hashedPassword, err := HashPassword(password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	result, err := IsMatch(usernameOrEmail, hashedPassword, db)
	if !result && err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/selectGame", http.StatusFound)
}

func BlindTestPage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/blindTest.html")
	temp.Execute(w, nil)
}

func DeafTestPage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/deafTest.html")
	temp.Execute(w, nil)
}

func PtitbacPage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/ptitBac.html")
	temp.Execute(w, nil)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/home.html")
	temp.Execute(w, nil)
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected successfully")

	reader(ws)
}
