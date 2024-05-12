package bdd

import (
	"GT/games"
	"database/sql"
)

func GetRoomID(roomName string) (int, error) {
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()
	var roomID int
	query := "SELECT id FROM ROOMS WHERE name = ?"
	err = db.QueryRow(query, roomName).Scan(&roomID)
	if err != nil {
		return 0, err
	}
	return roomID, nil
}

func GetRoomIDFromName(rommName string, db *sql.DB) (int, error) {
	var roomID int
	query := "SELECT id FROM ROOMS WHERE name = ?"
	err := db.QueryRow(query, rommName).Scan(&roomID)
	if err != nil {
		return 0, err
	}
	return roomID, nil
}

func QueryRoom(id int) (games.ROOM, error) {
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		return games.ROOM{}, err
	}
	defer db.Close()

	room := games.ROOM{}
	query := "SELECT created_by, max_player, name, id_game FROM ROOMS WHERE id = ?"
	err = db.QueryRow(query, id).Scan(&room.Created_by, &room.Max_players, &room.Name, &room.Id_game)
	if err != nil {
		return room, err
	}
	return room, nil
}
