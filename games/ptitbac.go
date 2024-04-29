package games

import (
	"log"
	"math/rand"
	"time"
)

func GenererLetters() string {
	letters := "abcdefghijklmnopqrstuvwxyz"
	rand.Seed(time.Now().UnixNano())
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




