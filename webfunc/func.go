package webfunc

import (
	. "GT/BDD"
	. "GT/Connect"
	"GT/games"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

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
	ArrayLetter  []string
	IsStarted    bool
	UsersInputs  map[string][]string
	Timer        int
	MaxRound     int
	CurrentRound int
	Categories   []string
	CurrentTime  int
	IsDone       bool
}

var arrayRoom = map[int]*PtitBacData{}

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
		jsonMsg := &struct {
			Id     int      `json:"id_room"`
			UserId int      `json:"id_user"`
			Done   bool     `json:"Done"`
			Data   []string `json:"data"`
			NextRound bool  `json: "NextRound"`
		}{}
		for {
			err := conn.ReadJSON(jsonMsg)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println(jsonMsg)
			room := arrayRoom[jsonMsg.Id]
			room.PtitBacConns.Add(conn)

			// fmt.Println("room connections :", room.PtitBacConns)
			defaultHandler := conn.CloseHandler()
			conn.SetCloseHandler(func(code int, text string) error {
				room.PtitBacConns.Delete(conn)
				return defaultHandler(code, text)
			})
			if jsonMsg.Done && !room.IsDone {
				fmt.Println(jsonMsg)
				room.IsStarted = true
				room.IsDone = true
				room.SendToRoom("end round")
			}
			if len(jsonMsg.Data) > 0 {
				username, err := QueryUserName(jsonMsg.UserId)
				if err != nil {
					log.Println("useridQuery : ", err)
				}
				room.UsersInputs[username] = jsonMsg.Data
				room.SendToRoom(map[string][]string{username: jsonMsg.Data})
				fmt.Println(room.UsersInputs)
			}
			if(jsonMsg.NextRound) {
				room.NextRound()
			}
		}
	case "loading":
		jsonMsg := &struct {
			Id     int
			UserId int `json:"id_user"`
			Start  bool
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

			r, err := QueryRoom(jsonMsg.Id)
			if err != nil {
				log.Println(err)
			}

			if jsonMsg.Start && r.Created_by == jsonMsg.UserId {
				room.IsStarted = true
				go room.StartTimer()
				room.SendToRoom("start game")
			} else {
				fmt.Println("non :", jsonMsg.Start, r.Created_by == jsonMsg.UserId)
			}
		}
	default:
		fmt.Println("Unknown game:", game)
	}
}

func (room *PtitBacData) SendToRoom(msg any) {
	for v := range room.PtitBacConns {
		if err := v.WriteJSON(msg); err != nil {
			log.Println(err)
			return
		}
	}
}

// MARK: nextRound
func (room *PtitBacData) NextRound() {
	if room.CurrentRound < room.MaxRound {
		room.CurrentRound++
		room.Letter = games.GenerateUniqueLetters(&room.ArrayLetter)
		room.SendToRoom("{letter : " + room.Letter + " }")
	} else {
		room.SendToRoom("end game")
	}
}

// MARK: Timer
func (room *PtitBacData) StartTimer() {
	go room.timerDecrease()
	time.Sleep(time.Duration(room.Timer) * time.Second)
	if !room.IsDone {
		room.SendToRoom("Timer expired!")
		log.Println(room.RoomLink, ": Timer expired!")
		room.IsDone = true
		room.SendToRoom("end round")
	}
}

func (room *PtitBacData) timerDecrease() {
	for room.CurrentTime != 0 {
		time.Sleep(time.Second)
		room.CurrentTime--
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

// MARK: CreateRoom
func CreateRoom(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	UserId := getUserIdFromPage(r)
	max_player, err := strconv.Atoi(r.FormValue("playersNumber"))
	if err != nil {
		log.Fatal(err)
	}

	timer, err := strconv.Atoi(r.FormValue("timerSeconds"))
	if err != nil {
		log.Fatal(err)
	}

	name := r.FormValue("name")
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	room := games.ROOM{
		Created_by:  UserId,
		Max_players: max_player,
		Name:        name,
		Id_game:     3,
	}
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	InsertRooms(room.Created_by, room.Max_players, room.Name, room.Id_game)
	roomID, err := GetRoomID(room.Name)
	if err != nil {
		log.Fatal(err)
	}

	maxRound, err := strconv.Atoi(r.FormValue("roundsNumber"))
	if err != nil {
		log.Fatal(err)
	}

	letters := []string{}
	letter := games.GenerateUniqueLetters(&letters)

	fmt.Println("catJSON :", r.FormValue("JSON"))
	var categories []string
	if err = json.Unmarshal([]byte(r.FormValue("JSON")), &categories); err != nil {
		log.Println(err)
	}
	fmt.Println(categories)

	//MARK: Init Room
	arrayRoom[roomID] = &PtitBacData{
		RoomLink:     "?room=" + fmt.Sprint(roomID),
		PtitBacConns: ConnSet{},
		Letter:       letter,
		ArrayLetter:  letters,
		IsStarted:    false,
		UsersInputs:  make(map[string][]string),
		Timer:        timer,
		MaxRound:     maxRound,
		CurrentRound: 1,
		Categories:   categories,
		CurrentTime:  timer,
	}

	http.Redirect(w, r, "/loadingPage?room="+strconv.Itoa(roomID), http.StatusFound)
}

func Loading(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/loading.html", "./template/websocket.html")
	r.ParseForm()
	roomId, err := strconv.Atoi(r.FormValue("room"))
	userId := getUserIdFromPage(r)
	db, err := sql.Open("sqlite3", "./BDD/table.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	query := "SELECT max_player FROM ROOMS WHERE id = ?"
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

	room, err := QueryRoom(roomId)
	if err != nil {
		log.Println(err)
	}

	temp.Execute(w, room.Created_by == userId)
}

func getUserIdFromPage(r *http.Request) int {
	cookie, err := r.Cookie("Id")
	if err != nil {
		log.Println(err)
	}
	userId, err := strconv.Atoi(cookie.Value)
	if err != nil {
		log.Println(err)
	}
	return userId
}

// MARK: PtitbacPage
func PtitbacPage(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("./pages/ptitBac.html", "./template/websocket.html")
	r.ParseForm()
	roomId, err := strconv.Atoi(r.FormValue("room"))
	if err != nil {
		log.Println(err)
	}

	roomRow, err := QueryRoom(roomId)
	if err != nil {
		log.Println(err)
	}
	userId := getUserIdFromPage(r)

	room := arrayRoom[roomId]
	temp.Execute(w, struct {
		Letter       string
		IsCreator    bool
		RoomId       int
		Categories   []string
		Time         int
		CurrentRound int
		MaxRound     int
		IsDone       bool
	}{
		Letter:       room.Letter,
		IsCreator:    roomRow.Created_by == userId,
		RoomId:       roomId,
		Categories:   room.Categories,
		Time:         room.CurrentTime,
		CurrentRound: room.CurrentRound,
		MaxRound:     room.MaxRound,
		IsDone:       room.IsDone,
	})
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
