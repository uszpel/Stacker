package main

import (
	"embed"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/* fonts/*
var assets embed.FS

const (
	myScreenWidth  = 600
	myScreenHeight = 800
	version        = "0.14.0"
)

func main() {
	g := &Game{}
	g.InitGame()
	ebiten.SetWindowSize(myScreenWidth, myScreenHeight)
	ebiten.SetWindowTitle("Stacker")
	fmt.Printf("Stacker version %v\n", version)
	fmt.Println("Copyright (c) 2025 uszpel.")

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
