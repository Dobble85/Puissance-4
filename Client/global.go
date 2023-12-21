package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

// Constantes définissant les paramètres généraux du programme.
const (
	globalWidth         = globalNumTilesX * globalTileSize
	globalHeight        = (globalNumTilesY + 1) * globalTileSize
	globalTileSize      = 100
	globalNumTilesX     = 7
	globalNumTilesY     = 6
	globalCircleMargin  = 5
	globalBlinkDuration = 60
	globalNumColorLine  = 3
	globalNumColorCol   = 3
	globalNumColor      = globalNumColorLine * globalNumColorCol
)

// Variables définissant les paramètres généraux du programme.
var (
	globalBackgroundColor    color.Color = color.NRGBA{R: 176, G: 196, B: 222, A: 255}
	globalGridColor          color.Color = color.NRGBA{R: 119, G: 136, B: 153, A: 255}
	globalTextColor          color.Color = color.NRGBA{R: 25, G: 25, B: 5, A: 255}
	globalSelectColor        color.Color = color.NRGBA{R: 0, G: 0, B: 255, A: 255}
	globalOponentSelectColor color.Color = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	globalBlackColor         color.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	smallFont                font.Face
	largeFont                font.Face
	globalTokenColors        [globalNumColor]color.Color = [globalNumColor]color.Color{
		color.NRGBA{R: 227, G: 52, B: 47, A: 255},
		color.NRGBA{R: 246, G: 153, B: 63, A: 255},
		color.NRGBA{R: 255, G: 237, B: 74, A: 255},
		color.NRGBA{R: 56, G: 193, B: 114, A: 255},
		color.NRGBA{R: 52, G: 144, B: 220, A: 255},
		color.NRGBA{R: 149, G: 97, B: 226, A: 255},
		color.NRGBA{R: 246, G: 109, B: 155, A: 255},
		color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	}
	offScreenImage *ebiten.Image
)
