package main

import (
	. "GT/BDD"
	. "GT/Connect"
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	CreateUserTable()
	CreateRoomsTable()
	CreateRoomUsersTable()
	CreateGamesTable()
	http.HandleFunc("/selectGame", func(w http.ResponseWriter, r *http.Request) {
		temp, _ := template.ParseFiles("./pages/selectGame.html")
		temp.Execute(w, nil)
	})
	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		temp, _ := template.ParseFiles("./pages/signin.html", "./template/signin.html")
		temp.Execute(w, nil)
	})
	http.HandleFunc("/createUser", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		temp, _ := template.ParseFiles("./pages/login.html")
		temp.Execute(w, nil)
	})
	http.HandleFunc("/loginUser", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/blindTest", func(w http.ResponseWriter, r *http.Request) {
		temp, _ := template.ParseFiles("./pages/blindtest.html")
		temp.Execute(w, nil)
	})

	http.HandleFunc("/deafTest", func(w http.ResponseWriter, r *http.Request) {
		temp, _ := template.ParseFiles("./pages/deafTest.html")
		temp.Execute(w, nil)
	})
	http.HandleFunc("/ptitbac", func(w http.ResponseWriter, r *http.Request) {
		temp, _ := template.ParseFiles("./pages/ptitBac.html")
		temp.Execute(w, nil)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		temp, _ := template.ParseFiles("./pages/home.html")
		temp.Execute(w, nil)
	})

	fs := http.FileServer(http.Dir("/static/"[1:]))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Listening on port 8080")
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
