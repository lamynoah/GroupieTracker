package bdd

import "database/sql"

func DeleteRoomsUser() error {
	db, err := sql.Open("sqlite3", "./BDD/table.db")
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