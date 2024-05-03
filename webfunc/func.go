package webfunc

import (
	. "GT/BDD"
	. "GT/Connect"
	"GT/games"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/websocket"
)

// type ptitBac struct {
// 	Artiste    string
// 	Album      string
// 	Groupe     string
// 	Instrument string
// 	Featuring  string
// }

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var ptitBacConns = ConnSet{}

type PtitBacData struct {
	// RoomName     string
	RoomLink     string
	PtitBacConns ConnSet
	Letter       string
	IsStarted    bool
	Inputs       map[string]games.Input
}

var arrayRoom = map[string]*PtitBacData{}

func reader(conn *websocket.Conn, game string) {
	ptitBacConns.Add(conn)
	if len(game) > 4 {
		game = game[4:]
	}
	switch game {
	case "blindTest":
		fmt.Println("game:", game)
	case "deafTest":
		fmt.Println("game:", game)
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			log.Println(string(p))

			if err := conn.WriteMessage(messageType, p); err != nil {
				log.Println(err)
				return
			}
		}
	case "ptitBac":
		ptitBacConns.Add(conn)
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			if string(p) == "Connected" {
				fmt.Println("je suis la", ptitBacConns)
				for v := range ptitBacConns {
					if err := v.WriteMessage(messageType, []byte("start game")); err != nil {
						log.Println(err)
						return
					}
				}
			}
			log.Println(string(p))

			if err := conn.WriteMessage(messageType, p); err != nil {
				log.Println(err)
				return
			}
		}
	case "loading":
		jsonMsg := &struct {
			Id    string
			Start bool
		}{}
		fmt.Println("all connections :", ptitBacConns)
		for {
			err := conn.ReadJSON(jsonMsg)
			if err != nil {
				log.Println(err)
				return
			}
			room := arrayRoom[jsonMsg.Id]
			room.PtitBacConns.Add(conn)
			defaultHandler := conn.CloseHandler()
			conn.SetCloseHandler(func(code int, text string) error {
				room.PtitBacConns.Delete(conn)
				return defaultHandler(code, text)
			})
			if jsonMsg.Start {
				room.IsStarted = true
				for v := range room.PtitBacConns {
					if err := v.WriteJSON("start game"); err != nil {
						log.Println(err)
						return
					}
				}
			}
		}
	default:
		fmt.Println("Unknown game:", game)
	}
}

func Select(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/selectGame.html")
	temp.Execute(w, nil)
}

func Lobby(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/lobby.html")
	temp.Execute(w, arrayRoom)
}

var dataError = struct {
	Error string
}{""}

func Signin(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/signin.html", "./template/signin.html")
	temp.Execute(w, dataError)
	dataError.Error = ""
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	UpperCase := regexp.MustCompile(`[A-Z]`)
	SpecialCar := regexp.MustCompile(`[^a-zA-Z0-9]`)
	hasNumber := regexp.MustCompile(`[0-9]`)
	if !UpperCase.MatchString(password) && !SpecialCar.MatchString(password) && !hasNumber.MatchString(password) && len(password) < 12 {
		dataError.Error = "Le mot de passe doit contenir au moins 1 nombre, 1 majuscule, 1 caractère spécial et doit être d'au moins 12 caractères"
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	if UserNameExist(username) {
		dataError.Error = "Cet Username est déja utilisé"
		http.Redirect(w, r, "/signin", http.StatusFound)
	}

	hasedPassword, err := HashPassword(password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = InsertUser(username, email, hasedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId, err := QueryUserId(username)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	UserCookies(w, userId)
	http.Redirect(w, r, "/selectGame", http.StatusFound)
}

func UserNameExist(username string) bool {
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := "SELECT COUNT(*) FROM Users WHERE username = ?"
	var count int
	err = db.QueryRow(query, username).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}

func Login(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/login.html")
	fmt.Println(Islogin(r))
	temp.Execute(w, dataError)
	dataError.Error = ""
}

func Connect(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	usernameOrEmail := r.FormValue("username")
	password := r.FormValue("password")

	result, err := IsMatch(usernameOrEmail, password, db)
	if !result || err != nil {
		fmt.Println(result, err)
		dataError.Error = "Username or password false"
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	userId, err := QueryUserId(usernameOrEmail)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	UserCookies(w, userId)
	http.Redirect(w, r, "/selectGame", http.StatusFound)
}

func BlindTestPage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/blindTest.html", "./template/websocket.html")
	temp.Execute(w, nil)
}

func DeafTestPage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/deafTest.html", "./template/websocket.html")
	temp.Execute(w, nil)
}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cookie, err := r.Cookie("Id")
	if err != nil {
		log.Fatal(err)
	}
	id, err := strconv.Atoi(cookie.Value)
	if err != nil {
		log.Fatal(err)
	}
	max_player, err := strconv.Atoi(r.FormValue("playersNumber"))
	if err != nil {
		log.Fatal(err)
	}
	room := games.ROOM{
		Created_by:  id,
		Max_players: max_player,
		Name:        r.FormValue("name"),
		Id_game:     3,
	}
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	InsertRooms(room.Created_by, room.Max_players, room.Name, room.Id_game)
	letters := []string{}
	letter = games.GenerateUniqueLetters(&letters)

	arrayRoom[room.Name] = &PtitBacData{
		"/loadingPage?room=" + room.Name,
		ConnSet{},
		letter,
		false,
		map[string]games.Input{},
	}
	http.Redirect(w, r, "/loadingPage?room="+room.Name, http.StatusFound)
}

func Result(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cookie, err := r.Cookie("Id")
	if err != nil {
		log.Println(err)
	}
	userName, err := QueryUserName(cookie.Value)
	if err != nil {
		log.Println(err)
	}
	// round := r.FormValue("roundsNumber")
	artiste := r.FormValue("artiste")
	album := r.FormValue("album")
	groupe := r.FormValue("groupe")
	instrument := r.FormValue("instrument")
	featuring := r.FormValue("featuring")
	input := games.Input{
		Artiste:    artiste,
		Album:      album,
		Groupe:     groupe,
		Instrument: instrument,
		Featuring:  featuring,
	}
	r.ParseForm()
	room := arrayRoom[r.FormValue("room")]
	room.Inputs[userName] = input
	// http.Redirect(w,r, "/")
}

func Loading(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/loading.html", "./template/websocket.html")
	r.ParseForm()
	roomId := r.FormValue("room")
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	query := "SELECT max_player FROM ROOMS WHERE name = ?"
	row := db.QueryRow(query, roomId)

	var maxPlayer int
	err = row.Scan(&maxPlayer)
	if err != nil {
		http.Redirect(w, r, "/lobby", http.StatusNotFound)
		return
	}
	if len(arrayRoom[roomId].PtitBacConns) > maxPlayer || arrayRoom[roomId].IsStarted {
		http.Redirect(w, r, "/lobby", http.StatusFound)
		return
	}

	temp.Execute(w, nil)
}

var letter string
var arrayInput []games.Input

func PtitbacPage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/ptitBac.html", "./template/websocket.html")
	r.ParseForm()
	temp.Execute(w, struct {
		Letter    string
		IsPlaying bool
	}{arrayRoom[r.FormValue("room")].Letter, true})
}

func SettingBacPage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/settingPagesPtitBac.html")

	temp.Execute(w, nil)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/home.html", "./template/websocket.html")
	temp.Execute(w, nil)
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected successfully")
	uri := r.RequestURI
	reader(ws, uri)
}
