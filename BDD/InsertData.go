package bdd

import ("database/sql"
_ "github.com/mattn/go-sqlite3")

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


