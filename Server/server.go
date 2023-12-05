package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const debug = true

type player struct {
	handler *bufio.ReadWriter
	ready   bool
	channel chan string
	id      int
}

type server struct {
	player1 *player
	player2 *player
	turn    int
}

func main() {
	// Init
	log.Println("[INFO] - Serveur démarré")
	server := server{
		player1: &player{channel: make(chan string), id: 1},
		player2: &player{channel: make(chan string), id: 2},
	}

	// Attente de connexion des joueurs
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("[ERROR] - Erreur lors de la création du listener:", err)
		return
	}

	log.Println("[INFO] - En attente de connexion")

	player1, err := listener.Accept()
	if err != nil {
		log.Fatalln("[ERROR] - Erreur lors de la connexion du joueur 1:", err)
		return
	}

	log.Println("[INFO] - Joueur 1 connecté")

	player2, err := listener.Accept()
	if err != nil {
		log.Fatalln("[ERROR] - Erreur lors de la connexion du joueur 2:", err)
		return
	}

	log.Println("[INFO] - Les deux joueurs sont connectés")
	server.player1.handler = bufio.NewReadWriter(bufio.NewReader(player1), bufio.NewWriter(player1))
	server.player2.handler = bufio.NewReadWriter(bufio.NewReader(player2), bufio.NewWriter(player2))
	defer listener.Close()

	go server.handlePlayer(1)
	go server.handlePlayer(2)

	// Choix des couleurs
	log.Println("[INFO] - En attente de la réponse des joueurs...")

	for server.player1.ready == false || server.player2.ready == false {
		time.Sleep(time.Millisecond * 100)
	}

	log.Println("[INFO] - Les deux joueurs ont choisi leur couleur")

	// Partie
	server.turn = 1
	for {
		println()
		log.Println("[INFO] - Début de la partie")
		server.broadcast("server:ready")
		server.player1.ready = false
		server.player2.ready = false

		for server.player1.ready == false || server.player2.ready == false {
			time.Sleep(time.Millisecond * 100)
		}

		log.Println("[INFO] - Partie terminée")
		println()
		// Fin de la partie

		log.Println("[INFO] - Synchronisation des joueurs...")

		server.player1.ready = false
		server.player2.ready = false

		for server.player1.ready == false || server.player2.ready == false {
			time.Sleep(time.Millisecond * 100)
		}

		log.Println("[INFO] -  Synchronisation de la partie")
		// TODO A revoir
		if server.turn == 1 {
			server.player1.send("0\n")
			server.player2.send("1\n")
		} else {
			server.player1.send("1\n")
			server.player2.send("0\n")
		}
		// TODO ------------------------------
		log.Println("[INFO] - Partie synchronisée")
		time.Sleep(time.Millisecond * 100)
	}

}

func (server *server) getPlayer(id int) *player {
	if id == 1 {
		return server.player1
	}
	return server.player2
}

func (player *player) send(msg string) {
	player.handler.WriteString(msg)
	player.handler.Flush()
	if debug {
		msg = strings.Replace(msg, "\n", "|", -1)
		log.Println("[SENT] - player "+fmt.Sprint(player.id)+" -> ", msg)
	}
}

func (player *player) receive() {
	for {
		msg, err := player.handler.ReadString('\n')
		if err != nil {
			log.Println("[ERROR] - Erreur lors de la lecture du message du serveur:", err)
			return
		}
		player.channel <- msg
		if debug {
			msg := strings.Replace(msg, "\n", "|", -1)
			log.Println("[RECEIVED] - player "+fmt.Sprint(player.id)+" -> ", msg)
		}
	}

}

func (server *server) broadcast(msg string) {
	server.player1.channel <- msg
	server.player2.channel <- msg
	if debug {
		msg = strings.Replace(msg, "\n", "|", -1)
		log.Println("[BROADCAST] - server -> ", msg)
	}
}

func (server *server) handlePlayer(id int) {
	player := server.getPlayer(id)
	other := server.getPlayer(3 - id)

	go player.receive()

	player.send(strconv.Itoa(id) + "\n")

	// Choix des couleurs
	colorChoice := false
	for {
		select {
		case msg := <-player.channel:
			if msg == "server:ready" {
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
	player.send("server:ready\n")

	// Partie + Synchro
	log.Println("[INFO] - Partie commencée - ", id)
	for {
		gameFinished := false
		for {
			// Partie
			select {
			case msg := <-player.channel:
				fmt.Println("")

				if msg == "server:game_finished" {
					player.ready = true
					gameFinished = true
					continue
				}
				if debug {
					log.Println("[DEBUG] - Partie -- ", id)
				}

				temp := strings.Split(msg, ", ")
				played := temp[0]
				gameFinished = temp[1] == "true\n"
				log.Println("[INFO] - Le joueur ", id, " a joué : ", played, " - Fin de partie : ", gameFinished)

				if gameFinished {
					player.ready = true
					other.channel <- "server:game_finished"
				}

				server.turn = 3 - id
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
				synchroFinished = msg == "server:ready"
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
