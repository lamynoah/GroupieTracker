package webfunc

import (
	"GT/bdd"
	"GT/connect"

	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

func reader(conn *websocket.Conn, game string) {
	// ptitBacConns.Add(conn)
	if len(game) > 4 {
		game = game[4:]
	}
	switch game {
	case "blindTest":
		fmt.Println("game:", game)
	case "deafTest":
		jsonMsg := &struct {
			Id        int    `json:"id_room"`
			UserId    int    `json:"id_user"`
			Done      bool   `json:"Done"`
			Input     string `json:"input"`
			NextRound bool   `json:"NextRound"`
			// Inputs    map[string]bool `json:"inputs"`
		}{}
		for {
			err := conn.ReadJSON(jsonMsg)
			if err != nil {
				log.Println(err)
				break
			}
			fmt.Println("msg :", jsonMsg)
			room := arrayRoomDeaftest[jsonMsg.Id]
			room.DeafTestConns.Store(conn, &sync.Mutex{})

			// fmt.Println("room connections :", room.PtitBacConns)
			defaultHandler := conn.CloseHandler()
			conn.SetCloseHandler(func(code int, text string) error {
				room.DeafTestConns.Delete(conn)
				return defaultHandler(code, text)
			})
			fmt.Println("done conditions :", jsonMsg.Done, room.IsDone)
			if jsonMsg.Done && !room.IsDone {
				fmt.Println("someone finished")
				room.IsStarted = true
				room.IsDone = true
				room.SendToRoom("end round")
			}
			switch {
			case jsonMsg.NextRound:
				room.NextRound()
			case len(jsonMsg.Input) > 0:
				room.UsersInputs.Store(jsonMsg.UserId, jsonMsg.Input)

				// case len(jsonMsg.Inputs) > 0:
				// 	fmt.Println("inputs :", jsonMsg.Inputs)
			}

		}
	case "ptitBac":
		jsonMsg := &struct {
			Id        int             `json:"id_room"`
			UserId    int             `json:"id_user"`
			Done      bool            `json:"Done"`
			Data      []string        `json:"data"`
			NextRound bool            `json:"NextRound"`
			Inputs    map[string]bool `json:"inputs"`
		}{}
		for {
			err := conn.ReadJSON(jsonMsg)
			if err != nil {
				log.Println(err)
				break
			}
			fmt.Println("msg :", jsonMsg)
			room := arrayRoom[jsonMsg.Id]
			room.PtitBacConns.Store(conn, &sync.Mutex{})

			// fmt.Println("room connections :", room.PtitBacConns)
			defaultHandler := conn.CloseHandler()
			conn.SetCloseHandler(func(code int, text string) error {
				room.PtitBacConns.Delete(conn)
				return defaultHandler(code, text)
			})
			fmt.Println("done conditions :", jsonMsg.Done, room.IsDone)
			if jsonMsg.Done && !room.IsDone {
				fmt.Println("someone finished")
				room.IsStarted = true
				room.IsDone = true
				room.SendToRoom("end round")
			}
			switch {
			case jsonMsg.NextRound:
				room.NextRound()
			case len(jsonMsg.Data) > 0:
				username, err := connect.QueryUserName(jsonMsg.UserId)
				if err != nil {
					log.Println("useridQuery : ", err)
				}
				room.UsersInputs.Store(username, jsonMsg.Data)
				ui := map[string][]string{}
				room.UsersInputs.Range(func(key, value any) bool {
					ui[key.(string)] = value.([]string)
					return true
				})
				room.SendToRoom(ui)
			case len(jsonMsg.Inputs) > 0:
				fmt.Println("inputs :", jsonMsg.Inputs)
				room.UsersPointsInputs = append(room.UsersPointsInputs, jsonMsg.Inputs)
				// AddScoreToPlayer(jsonMsg.Id, jsonMsg.UserId, 5)
				go room.CalcPoints()
			}

		}
	case "loading":
		jsonMsg := &struct {
			Id     int  `json:"id_room"`
			UserId int  `json:"id_user"`
			Start  bool `json:"start"`
		}{}
		for {
			err := conn.ReadJSON(jsonMsg)
			if err != nil {
				log.Println(err)
				break
			}
			room := arrayRoom[jsonMsg.Id]
			room.PtitBacConns.Store(conn, &sync.Mutex{})
			users, err := bdd.QueryRoomUsers(jsonMsg.Id)
			if err != nil {
				log.Println(err)
			}
			room.SendToRoom(&users)
			defaultHandler := conn.CloseHandler()
			conn.SetCloseHandler(func(code int, text string) error {
				room.PtitBacConns.Delete(conn)
				return defaultHandler(code, text)
			})

			r, err := bdd.QueryRoom(jsonMsg.Id)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(jsonMsg)
			if jsonMsg.Start && r.Created_by == jsonMsg.UserId {
				room.IsStarted = true
				go room.StartTimer()
				room.SendToRoom("start game")
			} else {
				fmt.Println("non :", jsonMsg.Start, r.Created_by == jsonMsg.UserId)
			}
		}
	case "loadingDeafTest":
		jsonMsg := &struct {
			Id     int  `json:"id_room"`
			UserId int  `json:"id_user"`
			Start  bool `json:"start"`
		}{}
		for {
			err := conn.ReadJSON(jsonMsg)
			if err != nil {
				log.Println(err)
				break
			}
			// fmt.Println("loadingDeafTest:", jsonMsg)
			room := arrayRoomDeaftest[jsonMsg.Id]
			room.DeafTestConns.Store(conn, &sync.Mutex{})
			users, err := bdd.QueryRoomUsers(jsonMsg.Id)
			if err != nil {
				log.Println(err)
			}
			room.SendToRoom(&users)
			defaultHandler := conn.CloseHandler()
			conn.SetCloseHandler(func(code int, text string) error {
				room.DeafTestConns.Delete(conn)
				return defaultHandler(code, text)
			})

			r, err := bdd.QueryRoom(jsonMsg.Id)
			if err != nil {
				log.Println(err)
			}
			if jsonMsg.Start && r.Created_by == jsonMsg.UserId {
				room.IsStarted = true
				go room.StartTimer()
				room.CurrentSong = GetSong()
				fmt.Println("before start :", room.CurrentSong)
				room.SendToRoom("start game")
			} else {
				fmt.Println("non :", jsonMsg.Start, r.Created_by == jsonMsg.UserId)
			}
		}
	default:
		fmt.Println("Unknown game:", game)
	}
}
