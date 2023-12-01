package main

import (
	"bufio"
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
		player1: &player{channel: make(chan string)},
		player2: &player{channel: make(chan string)},
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
		server.player1.ready = false
		server.player2.ready = false
		server.broadcast("server:ready")

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
		msg = strings.Replace(msg, "\n", "", -1)
		log.Println("[DEBUG] - Envoi du message au joueur : ", msg)
	}
}

func (player *player) receive() {
	msg, err := player.handler.ReadString('\n')
	if err != nil {
		log.Println("[ERROR] - Erreur lors de la lecture du message du serveur:", err)
		return
	}
	player.channel <- msg
	if debug {
		log.Print("[DEBUG] - Réception du message du joueur : ", msg)
	}
}

func (server *server) broadcast(msg string) {
	server.player1.channel <- msg
	server.player2.channel <- msg
}

func (server *server) handlePlayer(id int) {
	player := server.getPlayer(id)
	other := server.getPlayer(3 - id)

	player.send(strconv.Itoa(id) + "\n")

	// Choix des couleurs
	go player.receive()
	color_choice := false
	for {
		select {
		case msg := <-player.channel:
			if msg == "server:ready" {
				color_choice = true
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
				go player.receive()
			}
		default:
			// Do nothing
		}
		if color_choice {
			break
		}
	}
	player.send("server:ready\n")

	log.Println("[INFO] - Partie commencée - ", id)
	for {
		// Partie
		select {
		case msg := <-player.channel:
			if msg == "server:ready\n" {
				break
			}

			temp := strings.Split(msg, ", ")
			played := temp[0]
			game_finished := temp[1] == "true\n"
			log.Println("[INFO] - Le joueur ", id, " a joué : ", played, " - Fin de partie : ", game_finished)

			server.turn = 3 - id
			other.send(played + "\n")
			if game_finished {
				player.ready = true
				other.channel <- "server:ready\n"
				break
			}
			go player.receive()
		default:
			continue // Do nothing
		}

		// Synchro
		go player.receive()
		for {
			select {
			case msg := <-player.channel:
				if msg == "server:ready\n" {
					break
				}
				player.ready = true
			default:
				continue // Do nothing
			}
		}
	}
}
