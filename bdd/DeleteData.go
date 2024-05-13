package bdd

import "database/sql"

func DeleteRoomUsers() error {
	db, err := sql.Open("sqlite3", "./bdd/table.db")
	if err != nil {
		return err
	}
	defer db.Close()
	deleteQuery := "DELETE FROM ROOM_USERS"
	_, err = db.Exec(deleteQuery)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRooms() error {
	db, err := sql.Open("sqlite3", "./bdd/table.db")
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM ROOMS")
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM sqlite_sequence WHERE name=ROOMS")
	if err != nil {
		return err
	}
	return nil
}
