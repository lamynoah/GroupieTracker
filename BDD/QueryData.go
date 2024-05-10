package bdd

import "database/sql"

func GetRoomID(roomName string) (int, error) {
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()
	var roomID int
	query := "Select id FROM ROOMS WHERE name = ?"
	err = db.QueryRow(query, roomName).Scan(&roomID)
	if err != nil {
		return 0, err
	}
	return roomID, nil
}

func GetRoomIDFromName(rommName string, db *sql.DB) (int, error) {
	var roomID int
	query := "Select id FROM ROOMS WHERE name = ?"
	err := db.QueryRow(query, rommName).Scan(&roomID)
	if err != nil {
		return 0, err
	}
	return roomID, nil
}
