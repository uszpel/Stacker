package main

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

var idCounter = 0

type Vector struct {
	X float64
	Y float64
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
	sizeX := float64(b.Sprite.Bounds().Dx())
	sizeY := float64(b.Sprite.Bounds().Dy())
	return int(b.Position.X / sizeX), int(b.Position.Y / sizeY)
}

func (b *Block) Move(dX, dY float64) {
	if b.Moving {
		sizeX := float64(b.Sprite.Bounds().Dx())
		sizeY := float64(b.Sprite.Bounds().Dy())

		b.Position.X = b.Position.X + dX*sizeX
		b.Position.Y = b.Position.Y + dY*sizeY
	}
}

func (b *Block) Draw(screen *ebiten.Image) {
	sizeX := float64(b.Sprite.Bounds().Dx())
	sizeY := float64(b.Sprite.Bounds().Dy())

	for iy, y := range b.Shape {
		for ix := range y {
			if b.Shape[iy][ix] > 0 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(b.Position.X+float64(ix)*sizeX, b.Position.Y+float64(iy)*sizeY)
				screen.DrawImage(&b.Sprite, op)
			}
		}
	}
}

func NewBlock(img ebiten.Image, screenWidth float64, screenHeight float64) *Block {
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
