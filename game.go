package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"
	"os"

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
const stateReady = 5
const stateShowHighscores = 6
const stateExitRequested = 7

const datafile = "highscores.data"

type Game struct {
	CycleCounter int
	Direction    int
	Block        *Block
	Generator    BlockGenerator
	Board        [][]BoardEntry
	Score        int
	Lines        int
	FontSource   *text.GoTextFaceSource
	FontColor    color.Color
	State        int
	Level        int
	Menu         *ebiten.Image
	HighScore    *HighScore
}

type BoardEntry struct {
	Id     int
	Sprite int
}

func (g *Game) Update() error {
	g.checkKeyboardInput()
	if g.isRunning() {
		if g.Block == nil || !g.Block.Moving {
			completeLines := g.checkCompleteLines()
			if len(completeLines) > 0 {
				//log.Printf("Complete lines: %v", completeLines)
				g.removeCompleteLines(completeLines)
				g.Lines = g.Lines + len(completeLines)
				if g.Lines >= g.Level*10 {
					g.Level++
				}
			}
			if g.Block != nil {
				g.Score += g.Block.Score
			}
			g.Block = g.Generator.NewBlock(len(g.Board[0])/2-1, 2)
			if !g.checkBoard(*g.Block, 0, 1, false) {
				g.HighScore.InsertScore(*NewScore(g.Score, g.Lines, ""))
				//g.State = stateReadyToRestart
				g.initBoard()
				g.State = stateShowHighscores
				log.Printf("Game finished.")
				return nil
			}
			g.updateBoard(g.Block, g.Block.Id, g.Block.Sprite)
			g.Direction = dirNone
			//log.Printf("New block: %v\n", block.Id)
			//log.Print(g.printBoard())
		}

		if g.CycleCounter%20 == 0 {
			g.moveSideways(g.Block)
			if g.CycleCounter >= (60 / g.Level) {
				if g.checkState() {
					g.moveDown(g.Block)
				}
				g.CycleCounter = 0
			}
		}
		g.CycleCounter++
	}
	return nil
}

func (g Game) isRunning() bool {
	return g.State == stateRunningRequested || g.State == stateRunning || g.State == statePauseRequested
	//return g.State != statePaused && g.State != stateReadyToRestart && g.State != stateReady && g.State != stateShowHighscores
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
		distanceFromGround := g.calcDistanceFromGround(*block)
		g.updateBoard(block, 0, 0)

		if distanceFromGround == 0 {
			g.Direction = dirNone
		}
		if g.Direction == dirDown && g.checkBoard(*block, 0, distanceFromGround, false) {
			block.Move(0, distanceFromGround)
		} else {
			block.Move(0, 1)
		}
		g.Direction = dirNone
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

func (g *Game) checkState() bool {
	result := true
	switch g.State {
	case statePauseRequested:
		g.State = statePaused
		result = false
		//log.Print("Game paused.")
	case stateRunningRequested:
		g.State = stateRunning
		//log.Print("Game resumed.")
	}
	return result
}

func (g *Game) checkKeyboardInput() {
	switch g.State {
	case stateReady:
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			g.State = stateRunning
		}
		if ebiten.IsKeyPressed(ebiten.KeyH) {
			g.State = stateShowHighscores
		}
		if ebiten.IsKeyPressed(ebiten.KeyQ) {
			os.Exit(0)
		}
	case stateShowHighscores:
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			err := WriteHighscore(datafile, g.HighScore)
			if err != nil {
				log.Printf("Error writing highscore: %v", err)
			}
			g.State = stateReady
		}
	case stateRunning:
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
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			g.State = stateExitRequested
		}
	case statePaused:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.State = stateRunningRequested
		}
	case stateReadyToRestart:
		if ebiten.IsKeyPressed(ebiten.KeyY) {
			g.initBoard()
		}
	case stateExitRequested:
		if ebiten.IsKeyPressed(ebiten.KeyY) {
			g.State = stateReady
			g.initBoard()
		}
		if ebiten.IsKeyPressed(ebiten.KeyN) {
			g.State = stateRunning
		}
	}
}

func (g *Game) calcDistanceFromGround(block Block) int {
	result := len(g.Board) - 1
	gridX, gridY := block.getGridPosition()
	for iy, y := range block.Shape {
		for ix := range y {
			if block.Shape[iy][ix] > 0 {
				distance := 0
				for line := gridY; line+iy < len(g.Board)-1; line++ {
					if g.Board[line+iy][gridX+ix].Id == 0 {
						distance++
					}
					if g.Board[line+iy][gridX+ix].Id > 0 && g.Board[line+iy][gridX+ix].Id != block.Id {
						break
					}
				}
				if distance < result {
					result = distance
				}
			}
		}
	}
	return result
}

func (g *Game) InitGame() {
	g.CycleCounter = 0
	g.Generator.Init()
	g.mustLoadFont("fonts/arial_bold.ttf")
	g.HighScore, _ = ReadHighscore(datafile)
	g.initBoard()
}

func (g *Game) initBoard() {
	sizeX := myScreenWidth / g.Generator.Sprites[0].Bounds().Dx()
	sizeY := myScreenHeight / g.Generator.Sprites[0].Bounds().Dy()
	g.Board = make([][]BoardEntry, sizeY)
	for i := 0; i < sizeY; i++ {
		g.Board[i] = make([]BoardEntry, sizeX)
	}

	g.Block = nil
	g.Menu = nil
	g.FontColor = color.RGBA{0xcf, 0xcf, 0xcf, 0xff}
	g.Score = 0
	g.Lines = 0
	g.State = stateReady
	g.Level = 1
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
			if gridX < 0 || len(g.Board[0]) <= gridX+ix {
				result = false
				break
			} else if len(g.Board) <= gridY+iy+1 {
				block.Moving = false
				result = false
				break
			} else if block.Shape[iy][ix] > 0 && len(g.Board) > gridY+iy && len(g.Board[0]) > gridX+ix &&
				g.Board[gridY+iy][gridX+ix].Id > 0 {
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
		if curLine == len(line) {
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
	screen.Fill(color.RGBA{0x5f, 0x5f, 0x5f, 0xff})

	switch g.State {
	case stateReady:
		g.drawMenu(screen)
	case stateShowHighscores:
		g.drawHighscores(screen)
	default:
		g.drawBoard(screen)
	}
}

func (g *Game) drawBoard(screen *ebiten.Image) {
	playingField := ebiten.NewImage(myScreenWidth-25, myScreenHeight-80)
	playingField.Fill(color.RGBA{0x2f, 0x2f, 0x2f, 0xff})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(15, 63)
	screen.DrawImage(playingField, op)

	face := &text.GoTextFace{
		Source: g.FontSource,
		Size:   18,
	}
	g.prinText(screen, face, 15, 10, g.FontColor, fmt.Sprintf("Score: %d", g.Score))
	g.prinText(screen, face, 15, 35, g.FontColor, fmt.Sprintf("Lines: %d", g.Lines))
	g.prinText(screen, face, 515, 10, g.FontColor, fmt.Sprintf("Level: %2d", g.Level))

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

	boldFace := &text.GoTextFace{
		Source: g.FontSource,
		Size:   36,
	}
	switch g.State {
	case statePaused:
		fallthrough
	case statePauseRequested:
		g.prinText(screen, boldFace, 240, 380, color.White, "Paused")
	case stateReadyToRestart:
		g.prinText(screen, boldFace, 150, 380, color.White, "Game over. Restart? (y/n)")
	case stateExitRequested:
		g.prinText(screen, boldFace, 150, 380, color.White, "Exit game? (y/n)")
	}
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	if g.Menu == nil {
		g.Menu = ebiten.NewImage(myScreenWidth, myScreenHeight)
		for ix, b := range g.Board {
			for iy := range b {
				sprite := g.Generator.GetSprite(rand.IntN(1000) % len(g.Generator.Sprites))
				if rand.IntN(1000)%2 == 0 && (ix < 7 || ix > 13) && !(ix > 23) {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(15+float64(iy*sprite.Bounds().Dy()), 15+float64(ix*sprite.Bounds().Dx()))
					g.Menu.DrawImage(&sprite, op)
				}
			}
		}
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(g.Menu, op)

	boldFace := &text.GoTextFace{
		Source: g.FontSource,
		Size:   36,
	}
	g.prinText(screen, boldFace, 220, 280, g.FontColor, "Stacker")

	face := &text.GoTextFace{
		Source: g.FontSource,
		Size:   22,
	}
	g.prinText(screen, face, 250, 340, g.FontColor, "[S]tart")
	g.prinText(screen, face, 250, 370, g.FontColor, "[H]ighscores")
	g.prinText(screen, face, 250, 400, g.FontColor, "[Q]uit")
}

func (g *Game) drawHighscores(screen *ebiten.Image) {
	boldFace := &text.GoTextFace{
		Source: g.FontSource,
		Size:   36,
	}
	g.prinText(screen, boldFace, 120, 130, g.FontColor, "Stacker Highscores")

	face := &text.GoTextFace{
		Source: g.FontSource,
		Size:   22,
	}
	for index, score := range g.HighScore.Scores {
		g.prinText(screen, face, 150, 200+float64(index*30), g.FontColor,
			fmt.Sprintf("%2d. %d pts %d lines %s", index+1, score.Score, score.Lines, score.Name))
	}
}

func (g *Game) prinText(screen *ebiten.Image, face *text.GoTextFace, x float64, y float64, color color.Color, message string) {
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(x, y)
	opts.ColorScale.ScaleWithColor(color)
	text.Draw(screen, message, face, opts)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return myScreenWidth, myScreenHeight
}
