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
}

type fetchResponse struct {
	Lyrics string `json:"lyrics"`
}

func GetLyrics() (string, error) {
	songs := []Song{}
	content, err := os.ReadFile("./static/json/songs.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	err = json.Unmarshal(content, &songs)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	// fmt.Println(songs)

	index := rand.Intn(len(songs))
	selectedSong := songs[index]

	// trackTitleGTS = selectedSong.Title

	fmt.Println("Autheur: ", selectedSong.Author, " Nom de la musique: ", selectedSong.Title)
	// https://api.lyrics.ovh/v1/LCD Soundsystem/LCD Soundsystem
	url := fmt.Sprintf("https://api.lyrics.ovh/v1/"+selectedSong.Author+"/"+selectedSong.Title)
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var fetchResponse fetchResponse
	if err := json.NewDecoder(response.Body).Decode(&fetchResponse); err != nil {
		return "", err
	}

	return fetchResponse.Lyrics, nil
}
