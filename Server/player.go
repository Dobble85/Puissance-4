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
	_, err := player.handler.WriteString(msg)
	if err != nil {
		log.Println(Grey+"["+Red+"ERROR"+Grey+"]"+Reset+"- Erreur lors de l'envoi d'un message au joueur :", err)
		return
	}
	err = player.handler.Flush()
	if err != nil {
		log.Println(Grey+"["+Red+"ERROR"+Grey+"]"+Reset+"- Erreur lors de l'envoi d'un message au joueur :", err)
		return
	}
	if debug {
		log.Println(Grey+"["+Green+"SENT"+Grey+"]"+Reset+"- player <- ", strings.Replace(msg, "\n", "|", -1))
	}
}

func (player *player) receive(server *server) {
	for {
		msg, err := player.handler.ReadString('\n')
		if err != nil {
			log.Println(Grey + "[" + Cyan + "INFO" + Grey + "]" + Reset + "- Un joueur s'est déconnecté")
			for i, v := range server.players {
				if v == player {
					// Delete the player
					server.players = append(server.players[:i], server.players[i+1:]...)
					break
				}
			}
			for i, v := range server.games {
				if v.host == player || v.client == player {
					// Send a message to the other player
					// Delete the game
					server.games = append(server.games[:i], server.games[i+1:]...)
					if v.host == player {
						v.client.send("game:other_player_left\n")
					} else {
						v.host.send("game:other_player_left\n")
					}
					log.Println(Grey+"["+Cyan+"INFO"+Grey+"]"+Reset+"- Partie supprimée -", v.name, "("+fmt.Sprint(v.id)+")")
					break
				}
			}
			return
		}
		player.channel <- msg
		if debug {
			msg := strings.Replace(msg, "\n", "|", -1)
			log.Println(Grey+"["+Green+"RECEIVED"+Grey+"]"+Reset+"- player  -> ", msg)
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
				log.Println(Grey + "[" + Cyan + "INFO" + Grey + "]" + Reset + "- Partie créée - " + gameName + " (" + fmt.Sprint(game.id) + ")")
				return
			} else if temp[0] == "game:refresh" {
				availableGames := make([]string, len(server.games))
				for _, game := range server.games {
					if game.client == nil {
						for i := 0; i < len(availableGames); i++ {
							if availableGames[i] == "" {
								availableGames[i] = fmt.Sprint(game.id) + " - " + game.name
								break
							}
						}
					}
				}
				player.send(strings.Join(availableGames, ", ") + "\n")
			} else {
				log.Println(Grey+"["+Red+"ERROR"+Grey+"]"+Red+" player_handle "+Reset+"- Message inconnu:", msg)
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
