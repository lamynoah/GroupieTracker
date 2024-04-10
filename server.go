package main

import (
	"net/http"
)


 // local server
func main() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)
	http.ListenAndServe(":8080", nil)
}
