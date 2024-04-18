package games

import (
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
