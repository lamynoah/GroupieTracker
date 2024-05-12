package bdd

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func CreateUserTable() error {

	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		return err
	}
	defer db.Close()
	createTableQuery := `
        CREATE TABLE IF NOT EXISTS Users (
            UserID INTEGER PRIMARY KEY AUTOINCREMENT,
            Username TEXT NOT NULL UNIQUE,
            Password TEXT NOT NULL,
            Email TEXT NOT NULL UNIQUE
        );
    `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return err
	}
	return nil
}

func CreateRoomsTable() error {
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		return err
	}
	defer db.Close()

	createTableQuery := `
        CREATE TABLE IF NOT EXISTS ROOMS (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            created_by INTEGER NOT NULL,
            max_player INTEGER NOT NULL,
            name TEXT NOT NULL,
            id_game INTEGER,
            FOREIGN KEY (created_by) REFERENCES Users(UserID),
            FOREIGN KEY (id_game) REFERENCES GAMES(id)
        );
    `

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func CreateRoomUsersTable() error {
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		return err
	}
	defer db.Close()

	createTableQuery := `
        CREATE TABLE IF NOT EXISTS ROOM_USERS (
            id_room INTEGER,
            id_user INTEGER,
            score INTEGER,
            FOREIGN KEY (id_room) REFERENCES ROOMS(id),
            FOREIGN KEY (id_user) REFERENCES Users(UserID),
            PRIMARY KEY (id_room, id_user)
        );
    `

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func CreateGamesTable() error {
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		return err
	}
	defer db.Close()

	createTableQuery := `
        CREATE TABLE IF NOT EXISTS GAMES (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL
        );
    `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return err
	}
	return nil
}
