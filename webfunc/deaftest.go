package webfunc

import (
	"html/template"
	"math/rand"
	"net/http"
)

type Music struct {
	Title  string
	Lyrics string
}

func DeafTestPage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/deafTest.html", "./template/websocket.html")
	randomMusic := songs[(rand.Int() % len(songs))]
	temp.Execute(w, randomMusic)
}

// func (m Music) getLyrics() string { return m.Lyrics }

func (s Song) Check(formValue string) bool {
	if formValue == s.Title || formValue == s.Author {
		return true
	}
	return false
}
