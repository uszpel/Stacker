package main

import (
	"embed"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/* fonts/*
var assets embed.FS

const myScreenWidth = 600
const myScreenHeight = 800
const version = "0.1.0"

func main() {
	g := &Game{}
	g.CicleCounter = 0
	g.initBoard()
	ebiten.SetWindowSize(myScreenWidth, myScreenHeight)
	ebiten.SetWindowTitle("Stacker")
	fmt.Printf("Stacker version %v\n", version)
	fmt.Println("Copyright (c) 2025 uszpel.")

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
