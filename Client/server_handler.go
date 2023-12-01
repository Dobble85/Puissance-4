package main

import (
	"bufio"
	"log"
	"strings"
)

const debug = true

type server struct {
	handler *bufio.ReadWriter
	channel chan string
	ready   bool
	wait    bool
}

func (s *server) send(message string) {
	s.handler.WriteString(message)
	s.handler.Flush()
	if debug {
		log.Print("[DEBUG] - Message envoyé au serveur : ", message)
	}
}

func (s *server) receive() {
	response, _ := s.handler.ReadString('\n')
	s.channel <- strings.TrimSuffix(response, "\n")

	if debug {
		log.Print("[DEBUG] - Message reçu du serveur : ", response)
	}
}

func (s *server) waitUntilServerIsReady() {
	go s.receive()
	for {
		select {
		case <-s.channel:
			s.ready = true
			log.Println("Le serveur est prêt")
			return
		default:
			// Do nothing
		}
	}
}

func (g *game) getColor() {
	go g.server.receive()
	for {
		select {
		case message := <-g.server.channel:
				g.
			return
		default:
			// Do nothing
		}
	}
}
