package webfunc

import (
	"GT/games"
	"fmt"
	"log"
	"time"
)

// MARK: nextRound

func (room *PtitBacData) NextRound() {
	if room.CurrentRound < room.MaxRound {
		room.IsDone = false
		// fmt.Println("room isdone nextround :", room.IsDone)
		room.CurrentRound++
		room.Letter = games.GenerateUniqueLetters(&room.ArrayLetter)
		room.CurrentTime = room.Timer
		room.SendToRoom(struct {
			Letter string `json:"letter"`
			Time   int    `json:"time"`
		}{room.Letter, room.Timer})
		go room.StartTimer()
	} else if room.CurrentRound == room.MaxRound {
		// fmt.Println("room isdone nextround :", room.IsDone)
		room.SendToRoom("end game")
	}
}

// MARK: Timer

func (room *PtitBacData) StartTimer() {
	ended := room.timerDecrease(room.CurrentRound)
	fmt.Println("timer ended because of someone :", !ended)
	// time.Sleep(time.Duration(room.Timer) * time.Second)
	if !room.IsDone && ended {
		room.SendToRoom("end round")
		log.Println("round done because of someone or timeout")
		room.IsDone = true
	}
}

func (room *PtitBacData) timerDecrease(round int) bool {
	for room.CurrentTime != 0 {
		time.Sleep(time.Second)
		room.CurrentTime--
		if room.IsDone || (round < room.CurrentRound) {
			fmt.Println("timer stopped because of someone")
			return false
		}
	}
	return true
}
func (room *BlindTestData) NextRound() {
	if room.CurrentRound < room.MaxRounds {
		room.IsDone = false

		room.CurrentRound++
		room.CurrentTime = room.Timer
		room.SendToRoom(struct {
			// Letter string `json:"letter"`
			Time int `json:"time"`
		}{room.Timer})
		go room.StartTimer()
	} else if room.CurrentRound == room.MaxRounds {
		room.SendToRoom("end game")
	}
}

// MARK: Timer

func (room *BlindTestData) StartTimer() {
	ended := room.timerDecrease(room.CurrentRound)
	fmt.Println("timer ended because of someone :", !ended)
	// time.Sleep(time.Duration(room.Timer) * time.Second)
	if !room.IsDone && ended {
		room.SendToRoom("end round")
		log.Println("round done because of someone or timeout")
		room.IsDone = true
	}
}

func (room *BlindTestData) timerDecrease(round int) bool {
	for room.CurrentTime != 0 {
		time.Sleep(time.Second)
		room.CurrentTime--
		if room.IsDone || (round < room.CurrentRound) {
			fmt.Println("timer stopped because of someone")
			return false
		}
	}
	return true
}

func (room *DeafTestData) NextRound() {
	if room.CurrentRound < room.MaxRound {
		room.IsDone = false
		// fmt.Println("room isdone nextround :", room.IsDone)
		room.CurrentRound++
		// room.Letter = games.GenerateUniqueLetters(&room.ArrayLetter)
		room.CurrentTime = room.Timer
		room.SendToRoom(struct {
			Lyrics string `json:"lyrics"`
			Time   int    `json:"time"`
		}{room.CurrentSong.Lyrics, room.Timer})
		go room.StartTimer()
	} else if room.CurrentRound == room.MaxRound {
		// fmt.Println("room isdone nextround :", room.IsDone)
		room.SendToRoom("end game")
	}
}

// MARK: Timer

func (room *DeafTestData) StartTimer() {
	ended := room.timerDecrease(room.CurrentRound)
	fmt.Println("timer ended because of someone :", !ended)
	// time.Sleep(time.Duration(room.Timer) * time.Second)
	if !room.IsDone && ended {
		room.SendToRoom("end round")
		log.Println("round done because of someone or timeout")
		room.IsDone = true
	}
}

func (room *DeafTestData) timerDecrease(round int) bool {
	for room.CurrentTime != 0 {
		time.Sleep(time.Second)
		room.CurrentTime--
		if room.IsDone || (round < room.CurrentRound) {
			fmt.Println("timer stopped because of someone")
			return false
		}
	}
	return true
}
