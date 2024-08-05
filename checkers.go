package main

import (
	"log"
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// Loads icon from file.
func loadIcon() image.Image {
	file, err := os.Open("assets/icon.png")  // Image created with flaticon.com
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

func main() {
	game := NewGame()
	winIcon := loadIcon()

	ebiten.SetWindowSize(winWidth, winHeight)
	ebiten.SetWindowTitle("Checkers")
	ebiten.SetWindowIcon([]image.Image{winIcon})

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
