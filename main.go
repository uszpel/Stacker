package main

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var assets embed.FS

const myScreenWidth = 800
const myScreenHeight = 600

func main() {
	g := &Game{}
	g.CicleCounter = 0
	g.Blocks = make([]*Block, 0)
	g.initBoard()
	ebiten.SetWindowSize(myScreenWidth, myScreenHeight)
	ebiten.SetWindowTitle("Tetris")

	//g.printBoard()
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
