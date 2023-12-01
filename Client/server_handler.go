package main

import (
	"bufio"
	"log"
	"strconv"
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
	log.Println("[DEBUG] - getColor()")
	for {
		select {
		case message := <-g.server.channel:
			log.Println("[DEBUG] - getColor() : ", message)
			if message == "server:ready" {
				g.server.ready = true
				log.Println("[DEBUG] - Le serveur est prêt : getColor()")
				return
			} else {
				g.p2Color, _ = strconv.Atoi(message)
				go g.server.receive()
			}

		default:
			// Do nothing
		}
	}
}
