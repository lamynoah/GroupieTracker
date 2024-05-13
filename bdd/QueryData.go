package bdd

import (
	"GT/connect"
	"GT/games"
	"database/sql"
)

func GetRoomID(roomName string) (int, error) {
	db, err := sql.Open("sqlite3", "./bdd/table.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()
	var roomID int
	query := "SELECT id FROM ROOMS WHERE name = ?"
	err = db.QueryRow(query, roomName).Scan(&roomID)
	if err != nil {
		return 0, err
	}
	return roomID, nil
}

func GetRoomIDFromName(rommName string, db *sql.DB) (int, error) {
	var roomID int
	query := "SELECT id FROM ROOMS WHERE name = ?"
	err := db.QueryRow(query, rommName).Scan(&roomID)
	if err != nil {
		return 0, err
	}
	return roomID, nil
}

func QueryRoom(id int) (games.ROOM, error) {
	db, err := sql.Open("sqlite3", "./bdd/table.db")
	if err != nil {
		return games.ROOM{}, err
	}
	defer db.Close()

	room := games.ROOM{}
	query := "SELECT created_by, max_player, name, id_game FROM ROOMS WHERE id = ?"
	err = db.QueryRow(query, id).Scan(&room.Created_by, &room.Max_players, &room.Name, &room.Id_game)
	if err != nil {
		return room, err
	}
	return room, nil
}

func QueryRoomUsers(roomId int) ([]string, error) {
	db, err := sql.Open("sqlite3", "./bdd/table.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	users := []string{}
	query := "SELECT u.username FROM Users AS u INNER JOIN ROOM_USERS AS ru ON u.UserID = ru.id_user WHERE id_room = ?"
	rows, err := db.Query(query, roomId)
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var user string
		err := rows.Scan(&user)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}
	return users, nil
}

func QueryRoomUsersScores(roomId int) ([]games.ScoreBoard, error) {
	db, err := sql.Open("sqlite3", "./bdd/table.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	users := []games.ScoreBoard{}
	// for
	query := "SELECT id_user, score FROM ROOM_USERS WHERE id_room = ?"
	rows, err := db.Query(query, roomId)
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var user games.ScoreBoard
		var userId int
		err := rows.Scan(&userId, &user.Score)
		if err != nil {
			return users, err
		}
		user.Username, err = connect.QueryUserName(userId)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}
	return users, nil
}
