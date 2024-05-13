package webfunc

import (
	"GT/bdd"
	"GT/connect"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	UpperCase := regexp.MustCompile(`[A-Z]`)
	SpecialCar := regexp.MustCompile(`[^a-zA-Z0-9]`)
	hasNumber := regexp.MustCompile(`[0-9]`)
	if !UpperCase.MatchString(password) && !SpecialCar.MatchString(password) && !hasNumber.MatchString(password) && len(password) < 12 {
		dataError.Error = "Le mot de passe doit contenir au moins 1 nombre, 1 majuscule, 1 caractère spécial et doit être d'au moins 12 caractères"
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	if UserNameExist(username) {
		dataError.Error = "Cet Username est déja utilisé"
		http.Redirect(w, r, "/signin", http.StatusFound)
	}

	hasedPassword, err := connect.HashPassword(password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = bdd.InsertUser(username, email, hasedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId, err := connect.QueryUserId(username)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	connect.UserCookies(w, userId)
	http.Redirect(w, r, "/selectGame", http.StatusFound)
}

func UserNameExist(username string) bool {
	db, err := sql.Open("sqlite3", BDDPath)
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

func Connect(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", BDDPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	usernameOrEmail := r.FormValue("username")
	password := r.FormValue("password")

	result, err := connect.IsMatch(usernameOrEmail, password, db)
	if !result || err != nil {
		fmt.Println(result, err)
		dataError.Error = "Username or password false"
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	userId, err := connect.QueryUserId(usernameOrEmail)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	connect.UserCookies(w, userId)
	http.Redirect(w, r, "/selectGame", http.StatusFound)
}

func getUserIdFromPage(r *http.Request) int {
	cookie, err := r.Cookie("Id")
	if err != nil {
		log.Println(err)
	}
	userId, err := strconv.Atoi(cookie.Value)
	if err != nil {
		log.Println(err)
	}
	return userId
}

func getRoomIdFromPage(r *http.Request) int {
	r.ParseForm()
	room := r.FormValue("room")
	roomId, err := strconv.Atoi(room)
	if err != nil {
		log.Println(err)
	}
	return roomId
}
