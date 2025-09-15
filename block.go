package main

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

var idCounter = 0

type Vector struct {
	X int
	Y int
}

type Block struct {
	Id       int
	Shape    [][]int
	Sprite   ebiten.Image
	Position Vector
	Screen   Vector
	Moving   bool
}

func (b *Block) getGridPosition() (int, int) {
	return b.Position.X / b.Sprite.Bounds().Dx(), b.Position.Y / b.Sprite.Bounds().Dy()
}

func (b *Block) Move(dX, dY int) {
	if b.Moving {
		b.Position.X = b.Position.X + dX*b.Sprite.Bounds().Dx()
		b.Position.Y = b.Position.Y + dY*b.Sprite.Bounds().Dy()
	}
}

func (b *Block) Draw(screen *ebiten.Image) {
	for iy, y := range b.Shape {
		for ix := range y {
			if b.Shape[iy][ix] > 0 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(b.Position.X+ix*b.Sprite.Bounds().Dx()), float64(b.Position.Y+iy*b.Sprite.Bounds().Dy()))
				screen.DrawImage(&b.Sprite, op)
			}
		}
	}
}

func NewBlock(img ebiten.Image, screenWidth int, screenHeight int) *Block {
	block := &Block{}
	block.Id = newId()
	block.Sprite = img
	block.Position.X = 240
	block.Position.Y = 54
	block.Screen.X = screenWidth
	block.Screen.Y = screenHeight
	block.Shape = generateBlock()
	block.Moving = true
	return block
}

func newId() int {
	idCounter++
	return idCounter
}

func generateBlock() [][]int {
	return [][]int{
		{1, 1, 1},
		{1, 0, 0},
	}
}
