package main

import (
	"image"
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
	Sprite   int
	Position Vector
	Moving   bool
	Score    int
}

func (b *Block) getGridPosition() (int, int) {
	return b.Position.X, b.Position.Y
}

func (b *Block) Move(dX, dY int) {
	if b.Moving {
		b.Position.X = b.Position.X + dX
		b.Position.Y = b.Position.Y + dY
	}
}

func (b *Block) Rotate() {
	if b.Moving {
		newShape := make([][]int, len(b.Shape[0]))
		for i := range newShape {
			newShape[i] = make([]int, len(b.Shape))
		}

		for iy, y := range b.Shape {
			for ix := range y {
				newShape[ix][iy] = b.Shape[iy][len(y)-1-ix]
			}
		}
		b.Shape = newShape
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
		{1, 1, 1},
		{0, 0, 1},
	})
	b.Shapes = append(b.Shapes, [][]int{
		{1, 1},
		{1, 1},
	})
	b.Shapes = append(b.Shapes, [][]int{
		{1, 1, 1, 1},
	})
	b.Shapes = append(b.Shapes, [][]int{
		{1, 1, 0},
		{0, 1, 1},
	})
	b.Shapes = append(b.Shapes, [][]int{
		{0, 1, 1},
		{1, 1, 0},
	})
	b.Shapes = append(b.Shapes, [][]int{
		{1, 1, 1},
		{0, 1, 0},
	})

	b.Sprites = make([]ebiten.Image, 0)
	b.Sprites = append(b.Sprites, b.mustLoadImage("assets/blue_square.png"))
	b.Sprites = append(b.Sprites, b.mustLoadImage("assets/blue_square.png"))
	b.Sprites = append(b.Sprites, b.mustLoadImage("assets/green_square.png"))
	b.Sprites = append(b.Sprites, b.mustLoadImage("assets/red_square.png"))
	b.Sprites = append(b.Sprites, b.mustLoadImage("assets/yellow_square.png"))
	b.Sprites = append(b.Sprites, b.mustLoadImage("assets/yellow_square.png"))
	b.Sprites = append(b.Sprites, b.mustLoadImage("assets/purple_square.png"))
}

func (b BlockGenerator) NewBlock(x, y int) *Block {
	curColor := rand.IntN(1000) % len(b.Sprites)
	block := &Block{}
	block.Id = b.newId()
	block.Sprite = curColor
	block.Position.X = x
	block.Position.Y = y
	block.Shape = b.Shapes[curColor]
	block.Score = 5
	block.Moving = true
	return block
}

func (b BlockGenerator) GetSprite(index int) ebiten.Image {
	var result ebiten.Image
	if index < len(b.Sprites) {
		result = b.Sprites[index]
	}
	return result
}

func (b BlockGenerator) newId() int {
	idCounter++
	return idCounter
}

func (b *BlockGenerator) mustLoadImage(name string) ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return *ebiten.NewImageFromImage(img)
}
