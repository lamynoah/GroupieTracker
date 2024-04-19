package api

import (
	"GT/api"
	"GT/auth"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Track struct {
	Name    string `json:"name"`
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
}

func GetAccessToken() string {
	accessToken := auth.Token()
	return accessToken
}

func requestSong() {
	accessToken := api.GetAccessToken()
	playlistID := api.RequestPlaylist()
	trackID := selectRandomTrack(playlistID)
	url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", trackID)

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

	var trackData Track
	if err := json.NewDecoder(resp.Body).Decode(&trackData); err != nil {
		log.Fatal("Error decoding JSON", err)
	}

	fmt.Println("Track:", trackData.Name)
	fmt.Println("Artist:")
	for _, artist := range trackData.Artists {
		fmt.Println("-", artist.Name)
	}
}

func selectRandomTrack(playlistID string) string {
	tracks := extractTrackID(playlistID)

	if len(tracks) == 0 {
		log.Fatal("Playlist is empty")
	}

	rand.Seed(time.Now().UnixNano())

	randomIndex := rand.Intn(len(tracks))
	return tracks[randomIndex]
}

func extractTrackID(playlistID string) []string {
	accessToken := GetAccessToken()
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

	var playlistData struct {
		Tracks struct {
			Items []struct {
				Track struct {
					ID string `json:"id"`
				} `json:"track"`
			} `json:"items"`
		} `json:"tracks"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&playlistData); err != nil {
		log.Fatal("Error decoding JSON", err)
	}

	var trackID []string
	for _, item := range playlistData.Tracks.Items {
		trackID = append(trackID, item.Track.ID)
	}

	return trackID
}
