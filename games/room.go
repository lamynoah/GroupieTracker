package games

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

