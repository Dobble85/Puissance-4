package main

import (
	"log"
	"net"
)

type server struct {
	players  []*player
	games    []*game
	listener net.Listener
}

func (s server) handlePlayerConnection() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		player := newPlayer(conn)
		s.players = append(s.players, player)
		go player.receive(&s)
		go player.handle(&s)
		log.Println(Grey + "[" + Cyan + "INFO" + Grey + "]" + Reset + "- Nouveau joueur connecté")
	}
}

func (server *server) findGame(id int) int {
	for p, v := range server.games {
		if v.id == id {
			return p
		}
	}
	return -1
}
