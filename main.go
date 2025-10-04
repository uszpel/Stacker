package main

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/* fonts/*
var assets embed.FS

const myScreenWidth = 800
const myScreenHeight = 600

func main() {
	g := &Game{}
	g.CicleCounter = 0
	g.initBoard()
	ebiten.SetWindowSize(myScreenWidth, myScreenHeight)
	ebiten.SetWindowTitle("Stacker")

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
