package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

const debug = false

const (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
	Grey  = "\033[90m"
)

type server struct {
	handler *bufio.ReadWriter
	channel chan string
	ready   bool
	wait    bool
}

func (s *server) send(message string) {
	_, err := s.handler.WriteString(message)
	if err != nil {
		log.Println(Grey+"["+Red+"ERROR"+Grey+"]"+Reset+"- Erreur lors de l'envoi d'un message au serveur :", err)
		return
	}
	err = s.handler.Flush()
	if err != nil {
		log.Println(Grey+"["+Red+"ERROR"+Grey+"]"+Reset+"- Erreur lors de l'envoi d'un message au serveur :", err)
		return
	}
	if debug {
		log.Print(Grey+"["+Green+"SENT"+Grey+"]"+Reset+" - server <- ", message)
	}
}

func (s *server) receive() {
	for {
		response, err := s.handler.ReadString('\n')
		if err != nil {
			println("Connexion avec le serveur interrompue")
			// Close the game
			os.Exit(0)
			return
		}
		if response == "game:other_player_left\n" {
			println("Votre adversaire a quitté la partie")
			// Close the game
			os.Exit(0)
			return
		}
		s.channel <- strings.TrimSuffix(response, "\n")

		if debug {
			log.Print(Grey+"["+Green+"RECEIVED"+Grey+"]"+Reset+"- server -> ", response)
		}
	}
}

func (g *game) getColor() {
	log.Println(Grey + "[" + Green + "DEBUG" + Grey + "]" + Reset + "- getColor()")
	for {
		select {
		case message := <-g.server.channel:
			log.Println(Grey+"["+Green+"DEBUG"+Grey+"]"+Reset+"- getColor() : ", message)
			if message == "game:ready" {
				g.server.ready = true
				log.Println(Grey + "[" + Green + "DEBUG" + Grey + "]" + Reset + "- Le serveur est prêt : getColor()")
				return
			} else {
				g.p2Color, _ = strconv.Atoi(message)
			}
		default:
			// Do nothing
		}
	}
}
