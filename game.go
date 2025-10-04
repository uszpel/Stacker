package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const dirNone = 0
const dirLeft = 1
const dirRight = 2
const dirUp = 3
const dirDown = 4

const stateRunningRequested = 0
const stateRunning = 1
const stateReadyToRestart = 2
const statePauseRequested = 3
const statePaused = 4

type Game struct {
	CicleCounter int
	Direction    int
	Block        *Block
	Generator    BlockGenerator
	Board        [][]BoardEntry
	Score        int
	FontSource   *text.GoTextFaceSource
	State        int
}

type BoardEntry struct {
	Id     int
	Sprite int
}

func (g *Game) Update() error {
	g.checkKeyboardInput()
	if g.State != statePaused && g.State != stateReadyToRestart {
		if g.Block == nil || !g.Block.Moving {
			completeLines := g.checkCompleteLines()
			if len(completeLines) > 0 {
				log.Printf("Complete lines: %v", completeLines)
				g.removeCompleteLines(completeLines)
				g.Score = g.Score + len(completeLines)
			}

			g.Block = g.Generator.NewBlock(len(g.Board[0])/2-1, 0)
			if !g.checkBoard(*g.Block, 0, 1, false) {
				g.State = stateReadyToRestart
				log.Printf("Game finished.")
				return nil
			}

			//log.Printf("New block: %v\n", block.Id)
			//log.Print(g.printBoard())
		}

		if g.CicleCounter%20 == 0 {
			g.moveSideways(g.Block)
			if g.CicleCounter >= 60 {
				g.checkState()
				g.moveDown(g.Block)
				g.CicleCounter = 0
			}
		}
		g.CicleCounter++
	}
	return nil
}

func (g *Game) moveDown(block *Block) {
	if g.Direction == dirUp {
		if g.checkBoard(*block, 0, 0, true) {
			g.updateBoard(block, 0, 0)
			block.Rotate()
			g.updateBoard(block, block.Id, block.Sprite)
		}
		g.Direction = dirNone
	}
	if g.checkBoard(*block, 0, 1, false) {
		g.updateBoard(block, 0, 0)
		if g.Direction == dirDown && g.checkBoard(*block, 0, 3, false) {
			block.Move(0, 3)
			g.Direction = dirNone
		} else {
			block.Move(0, 1)
		}
		g.updateBoard(block, block.Id, block.Sprite)
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
			g.updateBoard(block, 0, 0)
			block.Move(step, 0)
			g.updateBoard(block, block.Id, block.Sprite)
		}
		g.Direction = dirNone
	}
}

func (g *Game) checkState() {
	switch g.State {
	case statePauseRequested:
		g.State = statePaused
		log.Print("Game paused.")
	case stateRunningRequested:
		g.State = stateRunning
		log.Print("Game resumed.")
	}
}

func (g *Game) checkKeyboardInput() {
	if g.State == stateRunning {
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			g.Direction = dirLeft
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			g.Direction = dirRight
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			g.Direction = dirUp
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			g.Direction = dirDown
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.State = statePauseRequested
		}
	} else {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			if g.State == statePaused {
				g.State = stateRunningRequested
			}
		}
	}
}

func (g *Game) initBoard() {
	g.Generator.Init()

	sizeX := myScreenWidth / g.Generator.Sprites[0].Bounds().Dx()
	sizeY := myScreenHeight / g.Generator.Sprites[0].Bounds().Dy()
	g.Board = make([][]BoardEntry, sizeY)
	for i := 0; i < sizeY; i++ {
		g.Board[i] = make([]BoardEntry, sizeX)
	}

	g.Block = nil
	g.Score = 0
	g.State = stateRunning
	g.mustLoadFont("fonts/arial.ttf")
}

func (g *Game) updateBoard(block *Block, id int, sprite int) {
	gridX, gridY := block.getGridPosition()
	if id > 0 {
		//log.Printf("Current grid position %v: (%v,%v)", block.Id, gridX, gridY)
	}
	for iy, y := range block.Shape {
		for ix := range y {
			if block.Shape[iy][ix] > 0 {
				g.Board[gridY+iy][gridX+ix] = BoardEntry{
					Id:     id,
					Sprite: sprite,
				}
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
				g.Board[gridY+iy][gridX+ix].Id > 0 {
				//log.Printf("Current grid check %v", g.Board[gridY+iy][gridX+ix])
				if g.Board[gridY+iy][gridX+ix].Id != g.Block.Id {
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
			if g.Board[ix][iy].Id != 0 {
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
		var newBoard [][]BoardEntry
		newBoard = append(newBoard, make([]BoardEntry, len(g.Board[0])))
		for lineIndex, line := range g.Board {
			if lineIndex != index {
				newBoard = append(newBoard, line)
			}
		}
		g.Board = newBoard
	}
}

func (g *Game) printBoard() string {
	result := "Board: \n"
	for ix, line := range g.Board {
		lineString := ""
		for iy := range line {
			lineString = fmt.Sprintf("%s%3d ", lineString, g.Board[ix][iy].Id)
		}
		result += lineString + "\n"
	}
	return result
}

func (g *Game) mustLoadFont(name string) {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}

	g.FontSource, err = text.NewGoTextFaceSource(f)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	face := &text.GoTextFace{
		Source: g.FontSource,
		Size:   18,
	}
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(10, 10)
	opts.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, fmt.Sprintf("Score: %v", g.Score), face, opts)

	for ix, b := range g.Board {
		for iy, e := range b {
			if e.Id != 0 {
				sprite := g.Generator.GetSprite(e.Sprite)
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(15+float64(iy*sprite.Bounds().Dy()), 15+float64(ix*sprite.Bounds().Dx()))
				screen.DrawImage(&sprite, op)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return myScreenWidth, myScreenHeight
}
