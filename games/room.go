package games

import (
	"database/sql"
)

type ROOM struct {
	Created_by  int
	Max_players int
	Name        string
	Id_game     int
}

type Data struct {
	id_game int
	ROOMS   []ROOM
}

func CreateNewRoom(db *sql.DB, room ROOM) error {
	query := `INSERT INTO ROOMS (created_by, max_players, name, id_game) Values (?,?,?) `
	_, err := db.Exec(query, room.Created_by, room.Max_players, room.Name, room.Id_game)
	if err != nil {
		return err
	}
	return nil
}
