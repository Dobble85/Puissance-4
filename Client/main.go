package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

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

	var conn net.Conn
	var err error
	for conn == nil {
		log.Println("[INFO] - Tentative de connexion au serveur")
		conn, err = net.Dial("tcp", ip)
		if err != nil {
			log.Println("[ERROR] - Echec de la connexion au serveur")
			time.Sleep(time.Second * 2)
		}
	}
	log.Println("[INFO] - Connexion au serveur réussie")
	fmt.Println("")
	defer conn.Close()

	log.Println("[INFO] - Je suis connecté")
	g.server.handler = bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	// Choix pour la création d'une partie ou la connexion à une partie existante

	for {
		go g.server.receive()
		println("Voulez-vous créer une partie ou en rejoindre une ?")
		println("1 - Créer une partie")
		println("2 - Rejoindre une partie")
		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			println("Veuillez entrer un choix valide\n\n\n\n\n")
			continue
		}
		if choice == 1 {
			// Création d'une partie
			var name string
			var password string
			for {
				print("Veuillez entrer le nom de la partie : ")
				_, err := fmt.Scanln(&name)
				if err != nil {
					println("Veuillez entrer un nom valide")
					continue
				}
				break
			} // Nom
			for {
				print("Veuillez entrer le mot de passe de la partie : ")
				_, err := fmt.Scanln(&password)
				if err != nil {
					println("Veuillez entrer un mot de passe valide")
					continue
				}
				break
			} // Mot de passe
			g.server.send("game:create, " + name + ", " + password + "\n")

			println("Partie créée avec succès")
			break

		} else if choice == 2 {
			// Récupération des parties disponibles
			go g.server.receive()
			g.server.send("game:refresh\n")
			msg := <-g.server.channel
			msg = strings.TrimSuffix(msg, "\n")
			games := strings.Split(msg, ", ")
			if games[0] == "" && len(games) == 1 {
				games = []string{}
			}
			// Rejoindre une partie
			println("Liste des parties disponibles :")
			if len(games) == 0 {
				println("Aucune partie disponible, veuillez en créer une !\n\n\n\n\n")
				continue
			}
			println("0 - Retour")
			for _, game := range games {
				println(game)
			}
			println()
			var gameId int
			var password string
			for {
				print("Veuillez entrer l'id de la partie à rejoindre : ")
				_, err := fmt.Scanln(&gameId)
				if err != nil {
					println("Veuillez entrer un id valide")
					continue
				}
				break
			} // Id
			if gameId == 0 {
				println("\n\n\n\n\n")
				continue
			}
			for {
				print("Veuillez entrer le mot de passe de la partie : ")
				_, err := fmt.Scanln(&password)
				if err != nil {
					println("Veuillez entrer un mot de passe valide")
					continue
				}
				break
			} // Mot de passe
			g.server.send("game:join, " + fmt.Sprint(gameId) + ", " + password + "\n")

			// On vérifie que le serveur accepte la connexion
			msg = <-g.server.channel
			if msg == "game:accepted" { // ?? "game:accepted"
				println("Connexion à la partie réussie")
				break
			} else if msg == "game:full" { // ?? "game:full"
				println("La partie est pleine\n\n\n\n\n")
			} else if msg == "_password" { // ?? "game:wrong_password"
				println("Le mot de passe est incorrect\n\n\n\n\n")
			} else if msg == "ound" { // ?? "game:not_found"
				println("La partie n'existe pas\n\n\n\n\n")
			} else {
				println("Une erreur est survenue\n\n\n\n\n")
			}
		} else {
			println("Veuillez entrer un choix valide\n\n\n\n\n")
			continue
		}
	}
	go g.server.receive()

	// Fin de l'ajout

	ebiten.SetWindowTitle("Puissance 4 - En attente de l'autre joueur")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(&g); err != nil { // Modifié
		log.Fatal(err)
	}

}
