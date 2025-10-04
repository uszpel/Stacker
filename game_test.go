package main

import (
	"testing"
)

func TestRemovingLines(t *testing.T) {
	game := &Game{}
	game.Board = [][]BoardEntry{
		{{0, 0}, {1, 1}, {0, 0}},
		{{0, 0}, {1, 1}, {0, 0}},
		{{1, 1}, {1, 1}, {2, 2}},
	}
	t.Log(game.printBoard())

	lines := game.checkCompleteLines()
	t.Logf("Lines: %v", lines)
	game.removeCompleteLines(lines)
	t.Log(game.printBoard())

	if !compareBoards(game.Board, [][]BoardEntry{
		{{0, 0}, {0, 0}, {0, 0}},
		{{0, 0}, {1, 1}, {0, 0}},
		{{0, 0}, {1, 1}, {0, 0}},
	}) {
		t.FailNow()
	}
}

func compareBoards(board1, board2 [][]BoardEntry) bool {
	result := len(board1) > 0 && len(board1) == len(board2) &&
		len(board1[0]) > 0 && len(board1[0]) == len(board2[0])
	if result {
		for ix, b := range board1 {
			for iy := range b {
				if board1[ix][iy].Id != board2[ix][iy].Id ||
					board1[ix][iy].Sprite != board2[ix][iy].Sprite {
					return false
				}
			}
		}
	}
	return result
}
