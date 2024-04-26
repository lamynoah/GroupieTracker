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

var musics = []Music{{"a", "aaa"}, {"b", "bbb"}, {"c", "ccc"}}

func DeafTestPage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/deafTest.html", "./template/websocket.html")
	randomMusic := musics[(rand.Int() % len(musics))]
	temp.Execute(w, randomMusic)
}

func (m Music) getLyrics() string { return m.Lyrics }

func (m Music) Check(formValue string) bool {
	if formValue == m.Title {
		return true
	}
	return false
}
