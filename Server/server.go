package main

import (
	"bufio"
	"log"
	"net"
	"strings"
	"time"
)

type player struct {
	handler  *bufio.ReadWriter
	response string
}

type server struct {
	listener net.Listener
	player1  *player
	player2  *player
}

func main() {
	// Init
	log.Println("[INFO] - Serveur démarré")
	server := server{player1: &player{}, player2: &player{}}

	// Attente de connexion des joueurs
	server.waitForPlayer()
	defer server.listener.Close()

	server.player1.send("ready\n")
	server.player2.send("ready\n")

	// Attente de la réponse du choix des couleurs
	log.Println("[INFO] - En attente de la réponse des joueurs...")
	go server.player1.receive()
	go server.player2.receive()

	for len(server.player1.response) == 0 || len(server.player2.response) == 0 {
		time.Sleep(time.Millisecond * 100)
	}

	log.Println("[INFO] - Les deux joueurs ont choisi leur couleur")
	server.player1.send("ready\n" + server.player2.response + "0\n")
	server.player2.send("ready\n" + server.player1.response + "1\n")

	// Partie
	turn := 1
	for {
		println()
		log.Println("[INFO] - Début de la partie")
		partie_finie := false
		for !partie_finie {
			partie_finie = server.playRound(turn)
			turn = 3 - turn
		}
		log.Println("[INFO] - Partie terminée")
		// Fin de la partie
		// Sync
		time.Sleep(time.Second * 20)
	}

}

func (server *server) playRound(playerId int) bool {
	player := server.getPlayer(playerId)
	other := server.getPlayer(3 - playerId)

	log.Println("[INFO] - Début du tour du joueur", playerId)

	// Attente du choix du joueur
	player.receive()

	// Traitement du choix du joueur
	// Case, partie_finie\n
	temp := strings.Split(player.response, ", ")
	case1 := temp[0]
	partie_finie := temp[1] == "true\n"

	// Envoi du choix du joueur à l'autre joueur
	other.send(case1 + "\n")

	return partie_finie
}

func (server *server) waitForPlayer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("[ERROR] - Erreur lors de la création du listener:", err)
		return
	}
	server.listener = listener

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
	//msg = strings.Replace(msg, "\n", "", -1)
	//log.Println("[DEBUG] - Envoi du message au joueur", msg)
}

func (player *player) receive() {
	player.response = ""
	msg, err := player.handler.ReadString('\n')
	if err != nil {
		log.Println("[ERROR] - Erreur lors de la lecture du message du serveur:", err)
		return
	}
	player.response = msg
	//log.Print("[DEBUG] - Réception du message du joueur", msg)
}
