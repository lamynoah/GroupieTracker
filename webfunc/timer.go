package webfunc

import (
	"GT/games"
	"log"
	"time"
)

// MARK: nextRound
func (room *PtitBacData) NextRound() {
	if room.CurrentRound < room.MaxRound {
		room.CurrentRound++
		room.Letter = games.GenerateUniqueLetters(&room.ArrayLetter)
		room.SendToRoom("{letter : " + room.Letter + " }")
	} else {
		room.SendToRoom("end game")
	}
}

// MARK: Timer
func (room *PtitBacData) StartTimer() {
	go room.timerDecrease()
	time.Sleep(time.Duration(room.Timer) * time.Second)
	if !room.IsDone {
		room.SendToRoom("Timer expired!")
		log.Println(room.RoomLink, ": Timer expired!")
		room.IsDone = true
		room.SendToRoom("end round")
	}
}

func (room *PtitBacData) timerDecrease() {
	for room.CurrentTime != 0 {
		time.Sleep(time.Second)
		room.CurrentTime--
	}
}
