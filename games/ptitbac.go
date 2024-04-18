package games

import (
	"fmt"
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

func StartTimer() {
	duration :=  5* time.Second
	time.Sleep(duration)
	fmt.Println("Timer expired!")
}

// func RecupInput(w http.ResponseWriter, r *http.Request) {
// 	artiste := r.FormValue("artiste")
// 	album := r.FormValue("album")
// 	groupe := r.FormValue("groupe")
// 	instrument := r.FormValue("instrument")
// 	featuring := r.FormValue("featuring")
// }
