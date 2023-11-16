package main

import (
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *game) waitForPlayerColorChoice() {
	msg := g.waitForServer()
	g.p2Color, _ = strconv.Atoi(msg)
	log.Println("Couleur du joueur 2:", g.p2Color)

	msg = g.waitForServer()
	g.turn, _ = strconv.Atoi(msg)
	g.wait = false
	log.Println("Tour du joueur:", g.turn)
	if g.turn != p1Turn {
		log.Println("Je suis le joueur 2")
		g.player_choice = make([]string, 0)
		ebiten.SetWindowTitle("Puissance 4 - En attente de l'autre joueur")
		go waitForPlayerChoice(*g.server_handler, &g.player_choice)
	} else {
		log.Println("Je suis le joueur 1")
		ebiten.SetWindowTitle("Puissance 4 - A toi de jouer !")
	}

	log.Println("Début de la partie")
}
