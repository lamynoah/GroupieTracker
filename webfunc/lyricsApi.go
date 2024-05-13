package webfunc

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	// "fmt"
	// "math/rand"
	// "net/http"
)

type Song struct {
	Author string
	Title  string
	Lyrics string
}

func GetSong() Song {
	songs := []Song{}
	content, err := os.ReadFile("./static/json/songs.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	err = json.Unmarshal(content, &songs)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	index := rand.Intn(len(songs))
	selectedSong := songs[index]

	fmt.Println("Autheur: ", selectedSong.Author, " Nom de la musique: ", selectedSong.Title)
	url := fmt.Sprintf("https://api.lyrics.ovh/v1/" + selectedSong.Author + "/" + selectedSong.Title)
	response, err := http.Get(url)
	if err != nil {
		return Song{}
	}
	defer response.Body.Close()

	var fetchResponse struct {
		Lyrics string `json:"lyrics"`
	}
	if err := json.NewDecoder(response.Body).Decode(&fetchResponse); err != nil {
		return Song{}
	}
	selectedSong.Lyrics = fetchResponse.Lyrics
	return selectedSong
}
