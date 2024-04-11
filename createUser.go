package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func createUserTable() error {

	db, err := sql.Open("sqlite3", "./BDD/Users.db")
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

func insertUser(username, email, password string) error {
	db, err := sql.Open("sqlite3", "./BDD/Users.db")
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

func hashPassword(password string) (string, error) {
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hasedPassword), nil
}
