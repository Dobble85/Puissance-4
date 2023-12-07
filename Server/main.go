package main

import (
	"log"
	"net"
	"time"
)

const debug = false

func main() {
	// Init
	log.Println(Grey + "[" + Cyan + "INFO" + Grey + "]" + Reset + "- Serveur démarré")
	server := &server{}

	// Attente de connexion des joueurs
	server.listener, _ = net.Listen("tcp", ":8080")
	defer server.listener.Close()
	go server.handlePlayerConnection() // Accepte les connexions des joueurs et les ajoute à la liste des joueurs

	for {
		time.Sleep(time.Hour * 24 * 365) // Boucle infinie, ça rompiche fort
	}
}
