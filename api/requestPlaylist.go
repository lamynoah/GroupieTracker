package main

import (
	"GT/auth"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Playlist struct {
	Name   string `json:"name"`
	Tracks struct {
		Items []struct {
			Track struct {
				Name    string `json:"name"`
				Artists []struct {
					Name string `json:"name"`
				} `json:"artists"`
			} `json:"track"`
		} `json:"items"`
	} `json:"tracks"`
}

func getAccessToken() string {
	accessToken := auth.Token()
	return accessToken
}

func main() {
	accessToken := getAccessToken()
	playlistID := "5Gu2ik0W12YoOexbFzYMTK" //playlist ID
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

	var playlistData Playlist
	if err := json.NewDecoder(resp.Body).Decode(&playlistData); err != nil {
		log.Fatal("Error decoding JSON", err)
	}

	fmt.Println("Playlist:", playlistData.Name)
	fmt.Println("Tracks:")
	for _, item := range playlistData.Tracks.Items {
		track := item.Track
		fmt.Println("Track:", track.Name)
		fmt.Println("Artists:")
		for _, artist := range track.Artists {
			fmt.Println("-", artist.Name)
		}
	}
}
