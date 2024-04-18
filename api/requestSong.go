package api

import (
	"GT/auth"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

//lint:ignore U1000 reason used later
func requestSong() {
	accessToken := GetAccessToken()
	trackID := "5wViaajeHHPZlEjBY9nhU3"
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
