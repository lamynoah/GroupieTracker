package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Token() string {
	authURL := "https://accounts.spotify.com/api/token"
	clientID := "17e0676008714fb5836169461b3e90f9"
	clientSecret := "c9ced88874d7417a9a1f78be532ff8df"

	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	requestBody := strings.NewReader("grant_type=client_credentials")

	req, err := http.NewRequest("POST", authURL, requestBody)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.Status)
		fmt.Println(string(body))
		return ""
	}

	var responseMap map[string]interface{}
	if err := json.Unmarshal(body, &responseMap); err != nil {
		fmt.Println("Error parsing response:", err)
		return ""
	}
	accessToken, ok := responseMap["access_token"].(string)
	if !ok {
		fmt.Println("Error: Access token not found in response")
		return ""
	}

	return accessToken
}
