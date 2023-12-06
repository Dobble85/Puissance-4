package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type player struct {
	handler  *bufio.ReadWriter
	ready    bool
	channel  chan string
	gameTurn int
}

func (player *player) send(msg string) {
	player.handler.WriteString(msg)
	player.handler.Flush()
	if debug {
		log.Println("[SENT] - player -> ", strings.Replace(msg, "\n", "|", -1))
	}
}

func (player *player) receive() {
	for {
		msg, err := player.handler.ReadString('\n')
		if err != nil {
			log.Println("[INFO] - Un joueur s'est déconnecté")
			return
		}
		player.channel <- msg
		if debug {
			msg := strings.Replace(msg, "\n", "|", -1)
			log.Println("[RECEIVED] - player "+fmt.Sprint(player.gameTurn)+" -> ", msg)
		}
	}

}

func (player *player) handle(server *server) {
	for {
		select {
		case msg := <-player.channel:
			msg = strings.Replace(msg, "\n", "", -1)
			temp := strings.Split(msg, ", ")
			if temp[0] == "game:join" {
				gameId, _ := strconv.Atoi(temp[1])
				password := temp[2]
				gameIndex := server.findGame(gameId) // Récupération de la partie
				if gameIndex == -1 {
					player.send("game:not_found\n")
					continue
				}
				game := server.games[gameIndex]

				if game.client == nil {
					if game.password == password {
						game.client = player
						player.gameTurn = 2
						player.send("game:accepted\n")
						log.Println("[INFO] - Le joueur 2 a rejoint la partie", game.id)
						game.start()
						return
					} else {
						player.send("game:wrong_password\n")
					}
				} else {
					player.send("game:full\n")
				}
			} else if temp[0] == "game:create" {
				gameName := temp[1]
				password := temp[2]
				game := &game{
					id:       len(server.games) + 1,
					name:     gameName,
					password: password,
					host:     player,
				}
				server.games = append(server.games, game)
				player.gameTurn = 1
				log.Println("[INFO] - Partie créée - " + gameName + " (" + fmt.Sprint(game.id) + ")")
			} else if temp[0] == "game:refresh" {
				availableGames := make([]string, len(server.games))
				for i, game := range server.games {
					if game.client == nil {
						availableGames[i] = fmt.Sprint(game.id) + " - " + game.name
					}
				}
				player.send(strings.Join(availableGames, ", ") + "\n")
			} else {
				log.Println("[ERROR] - Message inconnu:", msg)
			}
		default:
			// Do nothing
		}
	}
}

func newPlayer(conn net.Conn) *player {
	handler := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	return &player{
		handler: handler,
		channel: make(chan string),
	}
}
