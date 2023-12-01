package main

import (
	"bufio"
	"log"
	"net"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font/opentype"
)

// Mise en place des polices d'écritures utilisées pour l'affichage.
func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	smallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 30,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}

	largeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 50,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Création d'une image annexe pour l'affichage des résultats.
func init() {
	offScreenImage = ebiten.NewImage(globalWidth, globalHeight)
}

// Création, paramétrage et lancement du jeu.
func main() {
	g := game{server: &server{channel: make(chan string), ready: false}} // Modification
	ip := ""
	// Ajout de la connexion au serveur
	if len(os.Args) > 2 {
		log.Println("Usage:", os.Args[0], "ip")
		return
	} else if len(os.Args) == 1 {
		ip = "localhost:8080"
	} else {
		ip = os.Args[1]
	}

	log.Println("[INFO] - Je me connecte au serveur")
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Println("[ERROR] - Erreur lors de la connexion au serveur:", err)
		return
	}
	defer conn.Close()

	log.Println("[INFO] - Je suis connecté")
	g.server.handler = bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	go g.server.receive()

	// Fin de l'ajout

	ebiten.SetWindowTitle("Puissance 4 - En attente de l'autre joueur")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(&g); err != nil { // Modifié
		log.Fatal(err)
	}

}
