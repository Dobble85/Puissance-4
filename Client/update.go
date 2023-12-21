package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"strconv"
)

// Mise à jour de l'état du jeu en fonction des entrées au clavier.
func (g *game) Update() error {

	g.stateFrame++

	switch g.gameState {
	case titleState:
		if g.titleUpdate() && g.server.ready { // Modification
			g.server.ready = false //Ajout
			g.server.wait = false  // Ajout
			go g.getColor()        // Ajout
			ebiten.SetWindowTitle("Puissance 4 - Choisis ta couleur")
			g.gameState++
		}
	case colorSelectState:
		if g.colorSelectUpdate() && g.server.ready { // Modification
			g.server.ready = false //Ajout
			if g.turn == p1Turn {
				ebiten.SetWindowTitle("Puissance 4 - A toi de jouer !")
			} else {
				ebiten.SetWindowTitle("Puissance 4 - En attente de l'autre joueur")
			}
			g.gameState++
		}
	case playState:
		g.tokenPosUpdate()
		var lastXPositionPlayed int
		var lastYPositionPlayed int
		if g.turn == p1Turn {
			lastXPositionPlayed, lastYPositionPlayed = g.p1Update()
		} else {
			lastXPositionPlayed, lastYPositionPlayed = g.p2Update()
		}
		if lastXPositionPlayed >= 0 {
			finished, result := g.checkGameEnd(lastXPositionPlayed, lastYPositionPlayed)
			if finished {
				g.result = result
				if g.turn == p2Turn {
					// Ajout de l'envoi de la position du pion au serveur
					g.server.send(fmt.Sprint(lastXPositionPlayed) + ", true" + "\n")
				}
				g.server.wait = false
				g.server.ready = false
				ebiten.SetWindowTitle("Puissance 4 - Fin de partie")
				g.gameState++
			} else if g.turn == p2Turn {
				// Ajout de l'envoi de la position du pion au serveur
				g.server.send(fmt.Sprint(lastXPositionPlayed) + ", false" + "\n")
				ebiten.SetWindowTitle("Puissance 4 - En attente de l'autre joueur")
			} else {
				ebiten.SetWindowTitle("Puissance 4 - A toi de jouer !")
			}
		}
	case resultState:
		if g.resultUpdate() && g.server.ready { // Modification
			g.reset()

			if g.turn != p1Turn {
				ebiten.SetWindowTitle("Puissance 4 - En attente de l'autre joueur")
			} else {
				ebiten.SetWindowTitle("Puissance 4 - A toi de jouer !")
			}
			g.gameState = playState
		}
	}

	return nil
}

// Mise à jour de l'état du jeu à l'écran titre.
func (g *game) titleUpdate() bool {
	g.stateFrame = g.stateFrame % globalBlinkDuration

	if !g.server.ready {
		select {
		case message := <-g.server.channel:
			if message == "game:ready" {
				g.server.ready = true
				ebiten.SetWindowTitle("Puissance 4 - Appuyez sur entrée")
			}
		default:
			// Do nothing
		}
	}
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

// Mise à jour de l'état du jeu lors de la sélection des couleurs.
func (g *game) colorSelectUpdate() bool {

	changed := false
	col := g.p1Color % globalNumColorCol
	line := g.p1Color / globalNumColorLine

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		col = (col + 1) % globalNumColorCol
		if line*globalNumColorLine+col == g.p2Color {
			col = (col + 1) % globalNumColorCol
		}
		changed = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		col = (col - 1 + globalNumColorCol) % globalNumColorCol
		if line*globalNumColorLine+col == g.p2Color {
			col = (col - 1 + globalNumColorCol) % globalNumColorCol
		}
		changed = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		line = (line + 1) % globalNumColorLine
		if line*globalNumColorLine+col == g.p2Color {
			line = (line + 1) % globalNumColorLine
		}
		changed = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		line = (line - 1 + globalNumColorLine) % globalNumColorLine
		if line*globalNumColorLine+col == g.p2Color {
			line = (line - 1 + globalNumColorLine) % globalNumColorLine
		}
		changed = true
	}

	g.p1Color = line*globalNumColorLine + col

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.server.wait = true
		changed = true
	}

	if changed {
		g.server.send(fmt.Sprint(g.p1Color) + ", " + strconv.FormatBool(g.server.wait) + "\n")
	}
	return g.server.wait
}

// Gestion de la position du prochain pion à jouer par le joueur 1.
func (g *game) tokenPosUpdate() {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.tokenPosition = (g.tokenPosition - 1 + globalNumTilesX) % globalNumTilesX
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.tokenPosition = (g.tokenPosition + 1) % globalNumTilesX
	}
}

// Gestion du moment où le prochain pion est joué par le joueur 1.
func (g *game) p1Update() (int, int) {
	lastXPositionPlayed := -1
	lastYPositionPlayed := -1
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if updated, yPos := g.updateGrid(p1Token, g.tokenPosition); updated {
			g.turn = p2Turn
			lastXPositionPlayed = g.tokenPosition
			lastYPositionPlayed = yPos
		}
	}
	return lastXPositionPlayed, lastYPositionPlayed
}

// Gestion de la position du prochain pion joué par le joueur 2 et
// du moment où ce pion est joué.
func (g *game) p2Update() (int, int) {
	select {
	case message := <-g.server.channel:
		position, _ := strconv.Atoi(message)
		updated, yPos := g.updateGrid(p2Token, position)
		for ; !updated; updated, yPos = g.updateGrid(p2Token, position) {
			position = (position + 1) % globalNumTilesX
		}
		g.turn = p1Turn
		return position, yPos
	default:
		return -1, -1
	}
}

// Mise à jour de l'état du jeu à l'écran des résultats.
func (g game) resultUpdate() bool {
	select {
	case message := <-g.server.channel:
		g.turn, _ = strconv.Atoi(message)
		g.server.ready = true
		return g.server.wait
	default:
		if g.server.wait {
			return true
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.server.wait = true
			g.server.send("ready\n")
			ebiten.SetWindowTitle("Puissance 4 - Fin de partie - En attente de l'autre joueur")
			return true
		}
		return false
	}
}

// Mise à jour de la grille de jeu lorsqu'un pion est inséré dans la
// colonne de coordonnée (x) position.
func (g *game) updateGrid(token, position int) (updated bool, yPos int) {
	for y := globalNumTilesY - 1; y >= 0; y-- {
		if g.grid[position][y] == noToken {
			updated = true
			yPos = y
			g.grid[position][y] = token
			return
		}
	}
	return
}

// Vérification de la fin du jeu : est-ce que le dernier joueur qui
// a placé un pion gagne ? est-ce que la grille est remplie sans gagnant
// (égalité) ? ou est-ce que le jeu doit continuer ?
func (g game) checkGameEnd(xPos, yPos int) (finished bool, result int) {

	tokenType := g.grid[xPos][yPos]

	// horizontal
	count := 0
	for x := xPos; x < globalNumTilesX && g.grid[x][yPos] == tokenType; x++ {
		count++
	}
	for x := xPos - 1; x >= 0 && g.grid[x][yPos] == tokenType; x-- {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// vertical
	count = 0
	for y := yPos; y < globalNumTilesY && g.grid[xPos][y] == tokenType; y++ {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut gauche/bas droit
	count = 0
	for x, y := xPos, yPos; x < globalNumTilesX && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x+1, y+1 {
		count++
	}

	for x, y := xPos-1, yPos-1; x >= 0 && y >= 0 && g.grid[x][y] == tokenType; x, y = x-1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut droit/bas gauche
	count = 0
	for x, y := xPos, yPos; x >= 0 && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x-1, y+1 {
		count++
	}

	for x, y := xPos+1, yPos-1; x < globalNumTilesX && y >= 0 && g.grid[x][y] == tokenType; x, y = x+1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// egalité ?
	if yPos == 0 {
		for x := 0; x < globalNumTilesX; x++ {
			if g.grid[x][0] == noToken {
				return
			}
		}
		return true, equality
	}

	return
}
