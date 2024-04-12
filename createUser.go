package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserTable() error {

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

func InsertUser(username, email, password string) error {
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

func HashPassword(password string) (string, error) {
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hasedPassword), nil
}

func QueryNameUsers(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT Username FROM Users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		usernames = append(usernames, name)
	}
	return usernames, nil
}

func QueryEmailUsers(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT Email FROM Users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		Emails = append(Emails, email)
	}
	return Emails, nil
}

func QueryPasswordUsers(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT Password FROM Users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var Passwords []string
	for rows.Next() {
		var password string
		if err := rows.Scan(&password); err != nil {
			return nil, err
		}
		Passwords = append(Passwords, password)
	}
	return Passwords, nil
}

type User struct {
	email    string
	password string
}

func GetAllUsersDetails(db *sql.DB) (map[string]User, error) {
	detailsUsers := make(map[string]User)
	usernames, err := QueryNameUsers(db)
	if err != nil {
		return nil, err
	}
	emails, err := QueryEmailUsers(db)
	if err != nil {
		return nil, err
	}
	passwords, err := QueryPasswordUsers(db)
	if err != nil {
		return nil, err
	}
	for i, username := range usernames {
		detailsUsers[username] = User{
			emails[i],
			passwords[i],
		}
	}
	return detailsUsers, nil
}

func IsMatch(usernameOrEmail, password string, db *sql.DB) (bool, error) {
	users, err := GetAllUsersDetails(db)
	if err != nil {
		return false, err
	}
	for userName, userValue := range users {
		fmt.Println(userName, userValue)
		if usernameOrEmail == userName || usernameOrEmail == userValue.email {
			if password == userValue.password {
				return true, nil
			} else {
				return false, fmt.Errorf("Incorrect password: " + password + " != " + userValue.password)
			}
		}
	}
	return false, nil
}

// username or email => parcourir username = true

// passwor = 1235 => hash
// bdd.password => hash.password = 1235 et que user or email = same password [index]
