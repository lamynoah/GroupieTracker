package api

import (
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

type PlaylistData struct {
	Tracks struct {
		Items []struct {
			Track struct {
				ID string `json:"id"`
			} `json:"track"`
		} `json:"items"`
	} `json:"tracks"`
}

func GetAccessToken() string {
	return auth.Token()
}

func RequestSong() string {
	accessToken := GetAccessToken()
	playlistID := RequestPlaylist()
	trackID := selectRandomTrack(playlistID)
	trackData := getTrackData(accessToken, trackID)

	printTrackInfo(trackData)

	return trackID
}

func selectRandomTrack(playlistID string) string {
	tracks := extractTrackIDs(playlistID)

	if len(tracks) == 0 {
		log.Fatal("La playlist est vide")
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(tracks))
	return tracks[randomIndex]
}

func extractTrackIDs(playlistID string) []string {
	accessToken := GetAccessToken()
	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s", playlistID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Erreur lors de la création de la requête:", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Erreur lors de l'envoi de la requête", err)
	}
	defer resp.Body.Close()

	var playlistData PlaylistData
	if err := json.NewDecoder(resp.Body).Decode(&playlistData); err != nil {
		log.Fatal("Erreur lors du décodage du JSON", err)
	}

	var trackIDs []string
	for _, item := range playlistData.Tracks.Items {
		trackIDs = append(trackIDs, item.Track.ID)
	}

	return trackIDs
}

func getTrackData(accessToken, trackID string) Track {
	url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", trackID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Erreur lors de la création de la requête:", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Erreur lors de l'envoi de la requête", err)
	}
	defer resp.Body.Close()

	var trackData Track
	if err := json.NewDecoder(resp.Body).Decode(&trackData); err != nil {
		log.Fatal("Erreur lors du décodage du JSON", err)
	}

	return trackData
}

func printTrackInfo(trackData Track) {
	fmt.Println("Titre :", trackData.Name)
	fmt.Println("Artistes :")
	for _, artist := range trackData.Artists {
		fmt.Println("-", artist.Name)
	}
}
