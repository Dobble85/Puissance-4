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
		msg = strings.Replace(msg, "\n", "|", -1)
		log.Println("[SENT] - player "+fmt.Sprint(player.gameTurn)+" -> ", msg)
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
	availableGames := make([]string, len(server.games))
	for i, game := range server.games {
		if game.client == nil {
			availableGames[i] = fmt.Sprint(game.id) + " - " + game.name
		}
	}
	player.send(strings.Join(availableGames, ", ") + "\n")

	for {
		select {
		case msg := <-player.channel:
			temp := strings.Split(msg, ", ")
			if temp[0] == "game:join" {
				gameId, _ := strconv.Atoi(temp[1])
				password := temp[2]
				password = strings.TrimSuffix(password, "\n")
				game := server.games[gameId-1]
				if game.client == nil && game.password == password {
					game.client = player
					player.gameTurn = 2
					player.send("server:game_accepted\n")
					log.Println("[INFO] - Le joueur 2 a rejoint la partie", game.id)
					game.start()
					return
				} else {
					player.send("server:game_refused\n")
				}
			} else if temp[0] == "game:create" {
				gameName := temp[1]
				password := temp[2]
				password = strings.TrimSuffix(password, "\n")
				game := &game{
					id:       len(server.games) + 1,
					name:     gameName,
					password: password,
					host:     player,
				}
				server.games = append(server.games, game)
				player.gameTurn = 1
				player.send("server:game_accepted\n")
				log.Println("[INFO] - Partie créée - " + gameName + " (" + fmt.Sprint(game.id) + ")")
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
