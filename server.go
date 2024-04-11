package main

import (
	"net/http"
	"text/template"
)

func main() {
	createUserTable()
	http.HandleFunc("/signIn", func(w http.ResponseWriter, r *http.Request) {
		temp, _ := template.ParseFiles("./pages/signIn.html", "./template/signIn.html")
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		hasedPassword, err := hashPassword(password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = insertUser(username, email, hasedPassword)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		temp.Execute(w, nil)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		usernameOrEmail := r.FormValue("usernameOrEmail")
		password := r.FormValue("password")
		hasshedPassword, err := hashPassword(password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}




	})



	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)
	http.ListenAndServe(":8080", nil)
}
