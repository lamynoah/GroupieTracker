package webfunc

import (
	. "GT/BDD"
	. "GT/Connect"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
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

var dataError = struct {
	Error string
}{""}

func Signin(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/signin.html", "./template/signin.html")
	temp.Execute(w, dataError)
	dataError.Error = ""
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	// UpperCase := regexp.MustCompile(`[A-Z]`)
	// SpecialCar := regexp.MustCompile(`[^a-zA-Z0-9]`)
	// hasNumber := regexp.MustCompile(`[0-9]`)
	// if !UpperCase.MatchString(password) && !SpecialCar.MatchString(password) && !hasNumber.MatchString(password) && len(password) < 12 {
	// 	dataError.Error = "Le mot de passe doit contenir au moins 1 nombre, 1 majuscule, 1 caractère spécial et doit être d'au moins 12 caractères"
	// 	http.Redirect(w, r, "/signin", http.StatusFound)
	// 	return
	// }
	// if UserNameExist(username) {
	// 	dataError.Error = "Cet Username est déja utilisé"
	// 	http.Redirect(w, r, "/signin", http.StatusFound)
	// }

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

func UserNameExist(username string) bool {
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := "SELECT COUNT(*) FROM Users WHERE username = ?"
	var count int
	err = db.QueryRow(query, username).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}

func Login(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/login.html")
	temp.Execute(w, dataError)
	dataError.Error = ""
}

func Connect(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	usernameOrEmail := r.FormValue("username")
	password := r.FormValue("password")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	result, err := IsMatch(usernameOrEmail, password, db)
	if !result || err != nil {
		fmt.Println(result, err)
		dataError.Error = "Username or password false"
		http.Redirect(w, r, "/login", http.StatusFound)
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
