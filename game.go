package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const dirNone = 0
const dirLeft = 1
const dirRight = 2
const dirUp = 3
const dirDown = 4

type Game struct {
	CicleCounter int
	Direction    int
	Blocks       []*Block
	Generator    BlockGenerator
	Board        [][]int
}

func (g *Game) Update() error {
	block := g.getMovingBlock()
	if block == nil {
		completeLines := g.checkCompleteLines()
		if len(completeLines) > 0 {
			log.Printf("Complete lines: %v", completeLines)
			g.removeCompleteLines(completeLines)
		}

		block = g.Generator.NewBlock(len(g.Board[0])/2-1, 0)
		g.Blocks = append(g.Blocks, block)

		log.Printf("New block: %v\n", block.Id)
		//g.printBoard()
	}

	g.Direction = g.getDirection(g.Direction)
	if g.CicleCounter%20 == 0 {
		g.moveSideways(block)
		if g.CicleCounter >= 60 {
			g.moveDown(block)
			g.CicleCounter = 0
		}
	}
	g.CicleCounter++
	return nil
}

func (g *Game) moveDown(block *Block) {
	if g.Direction == dirUp {
		if g.checkBoard(*block, 0, 0, true) {
			g.updateBoard(block, 0)
			block.Rotate()
			g.updateBoard(block, block.Id)
		}
		g.Direction = dirNone
	}
	if g.checkBoard(*block, 0, 1, false) {
		g.updateBoard(block, 0)
		if g.Direction == dirDown && g.checkBoard(*block, 0, 3, false) {
			block.Move(0, 3)
			g.Direction = dirNone
		} else {
			block.Move(0, 1)
		}
		g.updateBoard(block, block.Id)
	} else {
		block.Moving = false
	}
}

func (g *Game) moveSideways(block *Block) {
	if g.Direction == dirLeft || g.Direction == dirRight {
		step := 1
		if g.Direction == dirLeft {
			step = -1
		}
		if g.checkBoard(*block, step, 0, false) {
			g.updateBoard(block, 0)
			block.Move(step, 0)
			g.updateBoard(block, block.Id)
		}
		g.Direction = dirNone
	}
}

func (g *Game) getDirection(curDirection int) int {
	result := curDirection
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		result = dirLeft
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		result = dirRight
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		result = dirUp
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		result = dirDown
	}
	return result
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

func (g *Game) getBlock(id int) *Block {
	var block *Block = nil
	for _, b := range g.Blocks {
		if b.Id == id {
			block = b
			break
		}
	}
	return block
}

func (g *Game) initBoard() {
	g.Generator.Init()

	sizeX := myScreenWidth / g.Generator.Sprites[0].Bounds().Dx()
	sizeY := myScreenHeight / g.Generator.Sprites[0].Bounds().Dy()
	g.Board = make([][]int, sizeY)
	for i := 0; i < sizeY; i++ {
		g.Board[i] = make([]int, sizeX)
	}
}

func (g *Game) updateBoard(block *Block, value int) {
	gridX, gridY := block.getGridPosition()
	if value > 0 {
		//log.Printf("Current grid position %v: (%v,%v)", block.Id, gridX, gridY)
	}
	for iy, y := range block.Shape {
		for ix := range y {
			if block.Shape[iy][ix] > 0 {
				g.Board[gridY+iy][gridX+ix] = value
			}
		}
	}
}

func (g *Game) checkBoard(block Block, dX int, dY int, rotate bool) bool {
	result := true
	gridX, gridY := block.getGridPosition()
	if rotate {
		block.Rotate()
	} else {
		gridX += dX
		gridY += dY
	}
	for iy, y := range block.Shape {
		for ix := range y {
			if gridX < 0 || len(g.Board[0]) <= gridX+ix+1 {
				result = false
				break
			} else if len(g.Board) <= gridY+iy {
				block.Moving = false
				result = false
				break
			} else if block.Shape[iy][ix] > 0 && len(g.Board) > gridY+iy && len(g.Board[0]) > gridX+ix &&
				g.Board[gridY+iy][gridX+ix] > 0 {
				//log.Printf("Current grid check %v: %v", block.Id, g.Board[gridY+iy][gridX+ix])
				if g.Board[gridY+iy][gridX+ix] != block.Id {
					result = false
					break
				}
			}
		}
	}
	return result
}

func (g *Game) checkCompleteLines() []int {
	completeLines := make([]int, 0)
	for ix, line := range g.Board {
		curLine := 0
		for iy := range line {
			if g.Board[ix][iy] != 0 {
				curLine++
			}
		}
		if curLine == len(line)-1 {
			completeLines = append(completeLines, ix)
		}
	}
	return completeLines
}

func (g *Game) removeCompleteLines(lines []int) {
	for _, index := range lines {
		var newBoard [][]int
		newBoard = append(newBoard, make([]int, len(g.Board[0])))
		for lineIndex, line := range g.Board {
			if lineIndex != index {
				newBoard = append(newBoard, line)
			}
		}
		g.Board = newBoard
	}
}

func (g *Game) printBoard() string {
	var result string
	for ix, line := range g.Board {
		lineString := ""
		for iy := range line {
			lineString = fmt.Sprintf("%s%3d ", lineString, g.Board[ix][iy])
		}
		result += lineString + "\n"
	}
	return result
}

func (g *Game) Draw(screen *ebiten.Image) {
	for ix, b := range g.Board {
		for iy, e := range b {
			if e != 0 {
				block := g.getBlock(e)
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(15+float64(iy*block.Sprite.Bounds().Dy()), 15+float64(ix*block.Sprite.Bounds().Dx()))
				screen.DrawImage(&block.Sprite, op)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return myScreenWidth, myScreenHeight
}
