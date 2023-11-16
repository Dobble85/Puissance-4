package main

import (
	"bufio"
	"log"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type server struct {
	handler  *bufio.ReadWriter
	ready    bool
	wait     bool
	response string
}

func (g *game) getColor() {
	g.server.receive()
	g.server.response = strings.TrimSuffix(g.server.response, "\n")
	temp := strings.Split(g.server.response, ", ")
	g.p2Color, _ = strconv.Atoi(temp[0])
	g.turn, _ = strconv.Atoi(temp[1])
	g.server.wait = false

	if g.turn != p1Turn {
		log.Println("Je suis le joueur 2")
		ebiten.SetWindowTitle("Puissance 4 - En attente de l'autre joueur")
		go g.server.receive()
	} else {
		log.Println("Je suis le joueur 1")
		ebiten.SetWindowTitle("Puissance 4 - A toi de jouer !")
	}

	log.Println("Début de la partie")
}

func (s *server) send(message string) {
	s.handler.WriteString(message)
	s.handler.Flush()
	//log.Print("[DEBUG] - Message envoyé au serveur : ", message)
}

func (s *server) receive() {
	s.response = ""
	s.response, _ = s.handler.ReadString('\n')
	//log.Print("[DEBUG] - Message reçu du serveur : ", s.response)
}

func (s *server) waitUntilServerIsReady() {
	s.receive()
	s.ready = true
	//log.Println("[DEBUG] - Serveur prêt")
}
