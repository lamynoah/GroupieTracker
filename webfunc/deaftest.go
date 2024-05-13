package webfunc

import (
	"html/template"
	"net/http"
)

type Music struct {
	Title  string
	Lyrics string
}

func DeafTestPage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/deaftest.html", "./template/websocket.html")
	r.ParseForm()
	roomId := getRoomIdFromPage(r)
	room := arrayRoomDeaftest[roomId]
	// fmt.Println(room.CurrentSong)
	temp.Execute(w, room)
}

// func (m Music) getLyrics() string { return m.Lyrics }

func (s Song) Check(formValue string) bool {
	if formValue == s.Title || formValue == s.Author {
		return true
	}
	return false
}
