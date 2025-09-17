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

func (b *Block) Rotate() {
	if b.Moving {
		newShape := make([][]int, len(b.Shape[0]))
		for i := range newShape {
			newShape[i] = make([]int, len(b.Shape))
		}

		for iy, y := range b.Shape {
			for ix := range y {
				newShape[ix][iy] = b.Shape[iy][ix]
			}
		}
		b.Shape = newShape
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

	b.Sprites = make([]ebiten.Image, 0)
	b.Sprites = append(b.Sprites, b.mustLoadImage("assets/blue_square.png"))
	b.Sprites = append(b.Sprites, b.mustLoadImage("assets/green_square.png"))
	b.Sprites = append(b.Sprites, b.mustLoadImage("assets/red_square.png"))
}

func (b BlockGenerator) NewBlock() *Block {
	curColor := rand.IntN(3)
	block := &Block{}
	block.Id = b.newId()
	block.Sprite = b.Sprites[curColor]
	block.Position.X = 240
	block.Position.Y = 54
	block.Shape = b.generateBlock(curColor)
	block.Moving = true
	return block
}

func (b BlockGenerator) newId() int {
	idCounter++
	return idCounter
}

func (b BlockGenerator) generateBlock(curColor int) [][]int {
	return b.Shapes[curColor]
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
