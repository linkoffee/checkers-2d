package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := NewGame()

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Checkers")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
