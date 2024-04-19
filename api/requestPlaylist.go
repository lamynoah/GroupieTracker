package api

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

//lint:ignore U1000 reason
func selectRandomPlaylist() string {
	playlistID := []string{
		"5Gu2ik0W12YoOexbFzYMTK",
		"37i9dQZF1DX188IBQOaMiA",
		"37i9dQZF1DXacPj7eARo6k",
		"3B8dwKgIMb0yNvvrEeqgyR",
	}

	rand.Seed((time.Now().UnixNano()))
	randomIndex := rand.Intn(len(playlistID))
	return playlistID[randomIndex]
}

//lint:ignore U1000 reason
func RequestPlaylist() string {
	accessToken := GetAccessToken()
	playlistID := selectRandomPlaylist() //playlist ID
	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s", playlistID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("error sending request", err)
	}
	defer resp.Body.Close()

	return playlistID
}
