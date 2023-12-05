package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type game struct {
	id       int
	name     string
	password string
	host     *player
	client   *player
	turn     int
}

func (g *game) start() {
	go g.handlePlayer(1)
	go g.handlePlayer(2)

	for g.host.ready == false || g.client.ready == false {
		time.Sleep(time.Millisecond * 100)
	}

	// Partie
	g.turn = 1
	for {
		println()
		log.Println("[INFO] - Début de la partie")
		g.broadcast("game:ready")
		g.host.ready = false
		g.client.ready = false

		for g.host.ready == false || g.client.ready == false {
			time.Sleep(time.Millisecond * 100)
		}

		g.host.ready = false
		g.client.ready = false

		for g.host.ready == false || g.client.ready == false {
			time.Sleep(time.Millisecond * 100)
		}

		if g.turn == 1 {
			g.host.send("0\n")
			g.client.send("1\n")
		} else {
			g.host.send("1\n")
			g.client.send("0\n")
		}

		log.Println("[INFO] - Partie synchronisée")
		time.Sleep(time.Millisecond * 100)
	}
}

func (g *game) getPlayer(id int) *player {
	if id == 1 {
		return g.host
	}
	return g.client
}

func (g *game) broadcast(msg string) {
	g.host.channel <- msg
	g.client.channel <- msg
	if debug {
		msg = strings.Replace(msg, "\n", "|", -1)
		log.Println("[BROADCAST] - game -> ", msg)
	}
}

func (g *game) handlePlayer(id int) {
	player := g.getPlayer(id)
	other := g.getPlayer(3 - id)

	player.send(strconv.Itoa(id) + "\n")

	// Choix des couleurs
	colorChoice := false
	for {
		select {
		case msg := <-player.channel:
			if msg == "game:ready" {
				colorChoice = true
			} else {
				temp := strings.Split(msg, ", ")

				ready := temp[1] == "true\n"
				if ready {
					if !player.ready {
						log.Println("[INFO] - Le joueur ", id, " est prêt")
					}
					player.ready = true
				}
				other.send(temp[0] + "\n")
			}
		default:
			// Do nothing
		}
		if colorChoice {
			break
		}
	}
	player.send("game:ready\n")

	// Partie + Synchro
	log.Println("[INFO] - Partie commencée - ", id)
	for {
		gameFinished := false
		for {
			// Partie
			select {
			case msg := <-player.channel:
				fmt.Println("")

				if msg == "game:game_finished" {
					player.ready = true
					gameFinished = true
					continue
				}

				temp := strings.Split(msg, ", ")
				played := temp[0]
				gameFinished = temp[1] == "true\n"

				if debug {
					log.Println("[DEBUG] - Joueur", id, "a joué", played)
				}
				if gameFinished {
					player.ready = true
					other.channel <- "game:game_finished"
				}

				g.turn = 3 - id
				other.send(played + "\n")
			default:
				// Do nothing
			}
			if gameFinished {
				break
			}
		}

		// Synchro
		synchroFinished := false
		for {
			select {
			case msg := <-player.channel:
				synchroFinished = msg == "game:ready"
				player.ready = true
			default:
				// Do nothing
			}
			if synchroFinished {
				break
			}
		}
		log.Println("[INFO] - Synchronisation de la partie - ", id)
	}
}
