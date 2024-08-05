package main

import (
	"bytes"
	"image/color"
	"io"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	winWidth   = 800
	winHeight  = 800
	boardRows  = 8
	boardCols  = 8
	cellSize   = winWidth / boardCols
	sampleRate = 44100
)

type Game struct {
	board         [][]*Piece
	pieces        []*Piece
	hoverX        int
	hoverY        int
	selectedPiece *Piece
	currentPlayer PieceType
	audioContext  *audio.Context
	bgMusicPlayer *audio.Player
	movePlayer    *audio.Player
	capturePlayer *audio.Player
	winPlayer     *audio.Player
	gameOver      bool
}

func NewGame() *Game {
	game := &Game{
		board:         make([][]*Piece, boardRows),
		pieces:        []*Piece{},
		currentPlayer: White,
	}

	for i := range game.board {
		game.board[i] = make([]*Piece, boardCols)
	}

	game.initPieces()
	game.initSounds()
	return game
}

func (g *Game) initPieces() {
	g.pieces = []*Piece{}
	for i := range g.board {
		for j := range g.board[i] {
			g.board[i][j] = nil
		}
	}

	// White checkers arrangement
	for y := 0; y < 3; y++ {
		for x := 0; x < boardCols; x++ {
			if (x+y)%2 != 0 {
				piece := NewPiece(White, x, y, "checkersObj/white_default.png", "checkersObj/white_king.png")
				g.pieces = append(g.pieces, piece)
				g.board[y][x] = piece
			}
		}
	}

	// Black checkers arrangement
	for y := boardRows - 3; y < boardRows; y++ {
		for x := 0; x < boardCols; x++ {
			if (x+y)%2 != 0 {
				piece := NewPiece(Black, x, y, "checkersObj/black_default.png", "checkersObj/black_king.png")
				g.pieces = append(g.pieces, piece)
				g.board[y][x] = piece
			}
		}
	}
}

func (g *Game) initSounds() {
	g.audioContext = audio.NewContext(sampleRate)

	moveSound, err := loadSound("sounds/move.wav")
	if err != nil {
		log.Fatal(err)
	}
	g.movePlayer = g.audioContext.NewPlayerFromBytes(moveSound)

	captureSound, err := loadSound("sounds/capture.wav")
	if err != nil {
		log.Fatal(err)
	}
	g.capturePlayer = g.audioContext.NewPlayerFromBytes(captureSound)

	winSound, err := loadSound("sounds/win.wav")
	if err != nil {
		log.Fatal(err)
	}
	g.winPlayer = g.audioContext.NewPlayerFromBytes(winSound)

	bgmSound, err := loadMP3Sound("sounds/bg_music.mp3")
	if err != nil {
		log.Fatal(err)
	}
	bgmStream := audio.NewInfiniteLoop(bytes.NewReader(bgmSound), int64(len(bgmSound)))
	g.bgMusicPlayer, err = g.audioContext.NewPlayer(bgmStream)
	if err != nil {
		log.Fatal(err)
	}
	g.bgMusicPlayer.Play()
}

func loadMP3Sound(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d, err := mp3.DecodeWithSampleRate(sampleRate, file)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, d); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func loadSound(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	d, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	buf := make([]byte, d.Length())
	if _, err := d.Read(buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func (g *Game) Update() error {
	if g.gameOver {
		if !g.winPlayer.IsPlaying() {
			g.resetGame()
		}
		return nil
	}

	if g.checkWin() {
		g.winPlayer.Rewind()
		g.winPlayer.Play()
		g.gameOver = true
		return nil
	}

	mx, my := ebiten.CursorPosition()
	g.hoverX = mx / cellSize
	g.hoverY = my / cellSize

	if g.hoverX >= boardCols || g.hoverY >= boardRows {
		return nil
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.selectedPiece == nil {
			g.selectedPiece = g.board[g.hoverY][g.hoverX]
			if g.selectedPiece == nil || g.selectedPiece.pieceType != g.currentPlayer {
				g.selectedPiece = nil
			}
		} else {
			if g.isValidMove(g.selectedPiece, g.hoverX, g.hoverY) {
				g.movePiece(g.selectedPiece, g.hoverX, g.hoverY)
				g.switchPlayer()
			} else {
				log.Println("Wrong Move")
			}
			g.selectedPiece = nil
		}
	}

	return nil
}

func (g *Game) isValidMove(piece *Piece, toX, toY int) bool {
	if toX < 0 || toY < 0 || toX >= boardCols || toY >= boardRows {
		return false
	}
	if g.board[toY][toX] != nil {
		return false
	}
	dx := abs(toX - piece.x)
	dy := abs(toY - piece.y)

	if piece.isKing {
		// The king can move in any direction
		return dx == dy && dx > 0 && (g.canKingMove(piece, toX, toY) || g.canKingCapture(piece, toX, toY))
	} else {
		// A regular checker can only move forward
		if piece.pieceType == White {
			if toY < piece.y {
				return false // White checker cannot move backwards
			}
		} else if piece.pieceType == Black {
			if toY > piece.y {
				return false // Black checker cannot move backwards
			}
		}

		if dx == 1 && dy == 1 {
			return true
		}
		if dx == 2 && dy == 2 {
			mx := (piece.x + toX) / 2
			my := (piece.y + toY) / 2
			if g.board[my][mx] != nil && g.board[my][mx].pieceType != piece.pieceType {
				return true
			}
		}
	}
	return false
}

func (g *Game) canKingMove(piece *Piece, toX, toY int) bool {
	mx, my := piece.x, piece.y
	stepX := (toX - piece.x) / abs(toX-piece.x)
	stepY := (toY - piece.y) / abs(toY-piece.y)
	enemyCount := 0
	for i := 1; i <= abs(toX-piece.x); i++ {
		mx += stepX
		my += stepY
		if g.board[my][mx] != nil {
			if g.board[my][mx].pieceType == piece.pieceType {
				return false
			} else {
				enemyCount++
			}
		}
	}
	return enemyCount <= 1
}

func (g *Game) canKingCapture(piece *Piece, toX, toY int) bool {
	dx := abs(toX - piece.x)
	dy := abs(toY - piece.y)
	if dx != dy {
		return false // The king must move diagonally
	}

	stepX := (toX - piece.x) / dx
	stepY := (toY - piece.y) / dy

	x, y := piece.x, piece.y
	enemyFound := false

	// Walk diagonally between current and target position
	for i := 1; i < dx; i++ {
		x += stepX
		y += stepY
		if g.board[y][x] != nil {
			if g.board[y][x].pieceType != piece.pieceType {
				if enemyFound {
					return false // More than one enemy checker found within the move
				}
				enemyFound = true
			} else {
				return false // Found your own checker on the way
			}
		}
	}

	// After finding an enemy checker, we check that the next cell is free
	if enemyFound {
		x += stepX
		y += stepY
		if x >= 0 && y >= 0 && x < boardCols && y < boardRows {
			return g.board[y][x] == nil // Next cell must be free
		}
	}

	return false
}

func (g *Game) movePiece(piece *Piece, toX, toY int) {
	dx := abs(toX - piece.x)
	dy := abs(toY - piece.y)
	if dx == 2 && dy == 2 {
		mx := (piece.x + toX) / 2
		my := (piece.y + toY) / 2
		g.board[my][mx] = nil
		for i, p := range g.pieces {
			if p.x == mx && p.y == my {
				g.pieces = append(g.pieces[:i], g.pieces[i+1:]...)
				break
			}
		}
		g.capturePlayer.Rewind()
		g.capturePlayer.Play()
	} else {
		g.movePlayer.Rewind()
		g.movePlayer.Play()
	}

	g.board[piece.y][piece.x] = nil
	piece.x = toX
	piece.y = toY
	g.board[toY][toX] = piece

	if (piece.pieceType == White && toY == boardRows-1) || (piece.pieceType == Black && toY == 0) {
		piece.isKing = true
		piece.updateImage()
	}
}

func (g *Game) switchPlayer() {
	if g.currentPlayer == White {
		g.currentPlayer = Black
	} else {
		g.currentPlayer = White
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y := 0; y < boardRows; y++ {
		for x := 0; x < boardCols; x++ {
			if (x+y)%2 == 0 {
				vector.DrawFilledRect(screen, float32(x*cellSize), float32(y*cellSize), float32(cellSize), float32(cellSize), color.RGBA{245, 245, 220, 255}, false)
			} else {
				vector.DrawFilledRect(screen, float32(x*cellSize), float32(y*cellSize), float32(cellSize), float32(cellSize), color.RGBA{175, 128, 79, 255}, false)
			}
		}
	}

	for _, piece := range g.pieces {
		piece.Draw(screen)
	}

	if g.selectedPiece != nil {
		g.selectedPiece.DrawHighlight(screen, color.RGBA{0, 255, 0, 100})
	}

	if g.gameOver {
		g.resetGame()
	}
}

func (g *Game) resetGame() {
	if !g.winPlayer.IsPlaying() {
		g.initPieces()
		g.currentPlayer = White
		g.gameOver = false
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return winWidth, winHeight
}

func (g *Game) checkWin() bool {
	whiteCount := 0
	blackCount := 0
	for _, piece := range g.pieces {
		if piece.pieceType == White {
			whiteCount++
		} else if piece.pieceType == Black {
			blackCount++
		}
	}

	// Check if either player has no pieces left
	if whiteCount == 0 {
		g.winPlayer.Rewind()
		g.winPlayer.Play()
		return true
	}
	if blackCount == 0 {
		g.winPlayer.Rewind()
		g.winPlayer.Play()
		return true
	}

	// Check if the game is a draw
	if !g.canAnyPieceMove() {
		g.winPlayer.Rewind()
		g.winPlayer.Play()
		return true
	}

	return false
}

func (g *Game) canAnyPieceMove() bool {
	for _, piece := range g.pieces {
		for y := 0; y < boardRows; y++ {
			for x := 0; x < boardCols; x++ {
				if g.isValidMove(piece, x, y) {
					return true
				}
			}
		}
	}
	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
