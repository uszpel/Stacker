package main

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const dirNone = 0
const dirLeft = 1
const dirRight = 2

type Game struct {
	CicleCounter int
	Direction    int
	Blocks       []*Block
	Board        [][]int
	Sprites      []ebiten.Image
}

func (g *Game) Update() error {
	block := g.getMovingBlock()
	if block == nil {
		block = NewBlock(g.Sprites[0], myScreenWidth, myScreenHeight)
		g.Blocks = append(g.Blocks, block)

		log.Printf("New block: %v\n", block.Id)
		g.printBoard()
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.Direction = dirLeft
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.Direction = dirRight
	}
	if g.CicleCounter%20 == 0 {
		g.updateBoard(block, 0)
		if g.Direction != dirNone {
			step := 1
			if g.Direction == dirLeft {
				step = -1
			}
			if g.checkBoard(block, step, 0) {
				block.Move(step, 0)
			}
			g.Direction = dirNone
		}
		if g.CicleCounter >= 60 {
			if g.checkBoard(block, 0, 1) {
				block.Move(0, 1)
			} else {
				block.Moving = false
			}
			g.CicleCounter = 0
			g.updateBoard(block, block.Id)
		}
	}
	g.CicleCounter++
	return nil
}

func (g *Game) getMovingBlock() *Block {
	var block *Block = nil
	for _, b := range g.Blocks {
		if b.Moving {
			block = b
			break
		}
	}
	return block
}

func (g *Game) initBoard() {
	g.Sprites = make([]ebiten.Image, 0)
	g.Sprites = append(g.Sprites, g.mustLoadImage("assets/blue_square.png"))

	sizeX := myScreenWidth / g.Sprites[0].Bounds().Dx()
	sizeY := myScreenHeight / g.Sprites[0].Bounds().Dy()
	g.Board = make([][]int, sizeY)
	for i := 0; i < sizeY; i++ {
		g.Board[i] = make([]int, sizeX)
	}
}

func (g *Game) updateBoard(block *Block, value int) {
	gridX, gridY := block.getGridPosition()
	if value > 0 {
		log.Printf("Current grid position %v: (%v,%v)", block.Id, gridX, gridY)
	}
	for iy, y := range block.Shape {
		for ix := range y {
			if block.Shape[iy][ix] > 0 {
				g.Board[gridY+iy][gridX+ix] = value
			}
		}
	}
}

func (g *Game) checkBoard(block *Block, dX, dY int) bool {
	result := true
	gridX, gridY := block.getGridPosition()
	gridX += dX
	gridY += dY
	for iy, y := range block.Shape {
		for ix := range y {
			if gridX < 0 || len(g.Board[0]) <= gridX+ix+1 {
				result = false
				break
			} else if len(g.Board) <= gridY+iy {
				block.Moving = false
				break
			} else if block.Shape[iy][ix] > 0 && len(g.Board) > gridY+iy && len(g.Board[0]) > gridX+ix &&
				g.Board[gridY+iy][gridX+ix] > 0 {
				log.Printf("Current grid check %v: %v", block.Id, g.Board[gridY+iy][gridX+ix])
				if g.Board[gridY+iy][gridX+ix] != block.Id {
					result = false
					break
				}
			}
		}
	}

	return result
}

func (g *Game) printBoard() {
	for ix, row := range g.Board {
		rowString := ""
		for iy := range row {
			rowString = fmt.Sprintf("%s%d ", rowString, g.Board[ix][iy])
		}
		fmt.Println(rowString)
	}
}

func (g *Game) mustLoadImage(name string) ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return *ebiten.NewImageFromImage(img)
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, b := range g.Blocks {
		b.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return myScreenWidth, myScreenHeight
}
