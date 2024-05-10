package games

import (
	"database/sql"
	"log"
	"math/rand"
	"time"
)

func GenererLetters() string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	index := rand.Intn(len(letters))
	return string(letters[index])
}

func UniqueLetter(arrayLetter []string, letter string) bool {
	for _, v := range arrayLetter {
		if v == letter {
			return false
		}
	}
	return true
}

func GenerateUniqueLetters(arrayLetter *[]string) string {
	var letter string
	for {
		letter = GenererLetters()
		if UniqueLetter(*arrayLetter, letter) {
			*arrayLetter = append(*arrayLetter, letter)
			return letter
		}
	}
}

func StartTimer(duration int) {
	time.Sleep(time.Duration(duration) * time.Second)
	log.Println("Timer expired!")
}

type Input struct {
	Artiste    string
	Album      string
	Groupe     string
	Instrument string
	Featuring  string
}

func IsElementUnique(arrayInput []Input) {
	mapArtiste := make(map[string]int)
	mapAlbum := make(map[string]int)
	mapGroupe := make(map[string]int)
	mapInstrument := make(map[string]int)
	mapFeaturing := make(map[string]int)
	for _, v := range arrayInput {
		mapArtiste[v.Artiste]++
		mapAlbum[v.Album]++
		mapGroupe[v.Groupe]++
		mapInstrument[v.Instrument]++
		mapFeaturing[v.Featuring]++
	}
	// Si réponse Unique score += 2 si réponse non unique score += 1
	// {"noah":2,"omar":1}
}

// artiste = noah && artiste = noah  == 1 points

type ScoreBoard struct {
	Username string
	Score    int
}

func ScoreBoardData(room int, db *sql.DB) []ScoreBoard {
	rows, err := db.Query("SELECT * FROM ROOM_USERS WHERE room_id=? ORDER BY score DES", room)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	dataScoreBoard := []ScoreBoard{}
	for rows.Next() {
		var username string
		var score int
		if err := rows.Scan(&username, &score); err != nil {
			log.Fatal(err)
		}
		dataScoreBoard = append(dataScoreBoard, ScoreBoard{Username: username, Score: score})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return dataScoreBoard
}
