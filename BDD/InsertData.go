package bdd

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InsertUser(username, email, password string) error {
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		return err
	}
	defer db.Close()

	insertQuery := `INSERT INTO USERS (Username, Password, Email) Values (?,?,?)`
	_, err = db.Exec(insertQuery, username, password, email)
	if err != nil {
		return err
	}
	return nil
}

func InsertRooms(created_by int, max_player int, name string, id_game int) error {
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		return err
	}
	defer db.Close()
	insertQuery := `INSERT INTO ROOMS (created_by, max_player, name, id_game) Values(?,?,?,?)`
	_, err = db.Exec(insertQuery, created_by, max_player, name, id_game)
	if err != nil {
		return err
	}
	return nil
}

func InsertRoomsUser(id_room, id_user, score int) error {
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		return err
	}
	defer db.Close()
	insertQuery := "INSERT INTO ROOM_USER(id_room, id_user, score) Values(?,?,?)"
	_, err = db.Exec(insertQuery, id_room, id_user, score)
	if err != nil {
		return err
	}
	return nil
}

func InsertGames(name string) error {
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		return err
	}
	defer db.Close()
	insertQuery := `INSERT INTO GAMES(name) Values (?)`
	_, err = db.Exec(insertQuery, name)
	if err != nil {
		return err
	}
	return nil
}
