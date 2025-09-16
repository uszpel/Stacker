package main

import (
	_ "image/png"
	"math/rand/v2"

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

type BlockGenerator struct {
	Sprites []ebiten.Image
	Shapes  [][][]int
}

func (b *BlockGenerator) Init() {
	b.Shapes = make([][][]int, 0)
	b.Shapes = append(b.Shapes, [][]int{
		{1, 1, 1},
		{1, 0, 0},
	})
	b.Shapes = append(b.Shapes, [][]int{
		{1, 1},
		{1, 1},
	})
	b.Shapes = append(b.Shapes, [][]int{
		{1, 1, 1},
	})
}

func (b BlockGenerator) NewBlock() *Block {
	block := &Block{}
	block.Id = b.newId()
	block.Sprite = b.Sprites[0]
	block.Position.X = 240
	block.Position.Y = 54
	block.Shape = b.generateBlock()
	block.Moving = true
	return block
}

func (b BlockGenerator) newId() int {
	idCounter++
	return idCounter
}

func (b BlockGenerator) generateBlock() [][]int {
	return b.Shapes[rand.IntN(3)]
}
