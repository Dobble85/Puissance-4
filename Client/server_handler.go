package main

import "bufio"

type server struct {
	handler  *bufio.ReadWriter
	ready    bool
	wait     bool
	response string
}

func (s *server) send(message string) {
	s.handler.WriteString(message)
	s.handler.Flush()
	//log.Print("[DEBUG] - Message envoyé au serveur : ", message)
}

func (s *server) receive() {
	s.response, _ = s.handler.ReadString('\n')
	//log.Print("[DEBUG] - Message reçu du serveur : ", s.response)
}

func (s *server) waitUntilServerIsReady() {
	s.receive()
	s.ready = true
	//log.Println("[DEBUG] - Serveur prêt")
}
