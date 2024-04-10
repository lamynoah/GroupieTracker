package main

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

func createUserTable() {
    
    db, err := sql.Open("sqlite3", "./BDD/Users.db")
    if err != nil {
        fmt.Println(err)
        return
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
        fmt.Println(err)
        return
    }
}