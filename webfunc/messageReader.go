package webfunc

import (
	"GT/connect"
	"GT/bdd"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

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
			Id        int      `json:"id_room"`
			UserId    int      `json:"id_user"`
			Done      bool     `json:"Done"`
			Data      []string `json:"data"`
			NextRound bool     `json:"NextRound"`
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
				username, err := connect.QueryUserName(jsonMsg.UserId)
				if err != nil {
					log.Println("useridQuery : ", err)
				}
				room.UsersInputs[username] = jsonMsg.Data
				room.SendToRoom(map[string][]string{username: jsonMsg.Data})
				fmt.Println(room.UsersInputs)
			}
			if jsonMsg.NextRound {
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

			r, err := bdd.QueryRoom(jsonMsg.Id)
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
