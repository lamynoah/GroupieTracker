package games

import (
	"database/sql"
)

type ROOM struct {
	Id          int
	Created_by  int
	Max_players int
	Name        string
}

type Data struct {
	TypeOfGame string
	ROOMS      []ROOM
}

func CreateNewRoom(db *sql.DB, room ROOM) error {
	query := `INSERT INTO ROOMS (created_by, max_players,name) Values (?,?,?) `

	_, err := db.Exec(query, room.Created_by, room.Max_players, room.Name)
	if err != nil {
		return err
	}
	return nil
}
