package main

import (
    "image/color"
    "log"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/hajimehoshi/ebiten/v2/vector"
)

type PieceType int

const (
    Empty PieceType = iota
    White
    Black
)

type Piece struct {
	pieceType PieceType
	x         int
	y         int
	isKing    bool
	image     *ebiten.Image
	kingImage *ebiten.Image
}

func NewPiece(pieceType PieceType, x, y int, imgPath string, kingImgPath string) *Piece {
    image, _, err := ebitenutil.NewImageFromFile(imgPath)
    if err != nil {
        log.Fatal(err)
    }
    kingImage, _, err := ebitenutil.NewImageFromFile(kingImgPath)
    if err != nil {
        log.Fatal(err)
    }
    return &Piece{
        pieceType: pieceType,
        x:         x,
        y:         y,
        image:     image,
        kingImage: kingImage,
    }
}

func (p *Piece) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.x*cellSize), float64(p.y*cellSize))

	if p.isKing {
		screen.DrawImage(p.kingImage, op)
	} else {
		screen.DrawImage(p.image, op)
	}
}

func (p *Piece) DrawHighlight(screen *ebiten.Image, clr color.Color) {
    vector.DrawFilledRect(screen, float32(p.x*cellSize), float32(p.y*cellSize), float32(cellSize), float32(cellSize), clr, false)
}

func (p *Piece) updateImage() {
    if p.isKing {
        p.image = p.kingImage
    }
}