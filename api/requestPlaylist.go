package api

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var playlistData struct {
	Tracks struct {
		Items []struct {
			Track struct {
				ID string `json:"id"`
			} `json:"track"`
		} `json:"items"`
	} `json:"tracks"`
}

//lint:ignore U1000 reason
func RequestPlaylist() string {
	accessToken := GetAccessToken()
	playlistID := selectRandomPlaylist()
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

//lint:ignore U1000 reason
func selectRandomPlaylist() string {
	playlistID := []string{
		"16KcsJ0bML5NoxCptRd3SK",
		"0FaEfQoznBqCPbhlUFFNlO",
	}

	rand.Seed((time.Now().UnixNano()))
	randomIndex := rand.Intn(len(playlistID))
	return playlistID[randomIndex]
}
