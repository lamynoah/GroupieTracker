package main

import (
	"database/sql"
	"net/http"
	"text/template"
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
	http.HandleFunc("/signIn", func(w http.ResponseWriter, r *http.Request) {
		temp, _ := template.ParseFiles("./pages/signIn.html", "./template/signIn.html")
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

	http.HandleFunc("/logIn", func(w http.ResponseWriter, r *http.Request) {
		temp, _ := template.ParseFiles("./pages/logIn.html")
		temp.Execute(w, nil)
	})
	http.HandleFunc("/loginUser", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("sqlite3", "./BDD/Users.db")
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

	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)
	http.ListenAndServe(":8080", nil)
}
