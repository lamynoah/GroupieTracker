package main

import (
	"fmt"
	"net/http"
)

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Callback received!")
}

func main() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	http.HandleFunc("/callback", CallbackHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("[ERROR] - Server could not start properly.\n ", err)
	}
}
