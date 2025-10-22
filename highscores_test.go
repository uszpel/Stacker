package main

import (
	"os"
	"testing"
)

const datafile = "highscores.test.data"

func TestReadHighscoreWithoutFile(t *testing.T) {
	os.Remove(datafile)
	score, err := readHighscore(datafile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(score.Scores) > 0 {
		t.Errorf("Scores not empty as expected")
	}
}

func TestReadHighscoreWithFile(t *testing.T) {
	writeHighscore(datafile, createTestScore())
	defer os.Remove(datafile)

	score, err := readHighscore(datafile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(score.Scores) != 1 {
		t.Errorf("Scores has more than one entry")
	}
}

func createTestScore() *HighScore {
	testscore := &HighScore{
		Scores: make([]Score, 0),
	}
	testscore.Scores = append(testscore.Scores, Score{Rank: 1, Score: 100, Name: "Test"})
	return testscore
}
