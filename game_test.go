package main

import "testing"

func TestRemovingLines(t *testing.T) {
	game := &Game{}
	game.Board = [][]int{
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{1, 1, 2, 0},
	}
	game.printBoard()

	lines := game.checkCompleteLines()
	t.Logf("Lines: %v", lines)
	game.removeCompleteLines(lines)
	game.printBoard()
}
