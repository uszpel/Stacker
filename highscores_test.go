package main

import (
	"os"
	"testing"
)

const datafile = "highscores.test.data"

func TestReadHighscoreWithoutFile(t *testing.T) {
	os.Remove(datafile)
	score, err := ReadHighscore(datafile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(score.Scores) > 0 {
		t.Errorf("Scores not empty as expected")
	}
}

func TestReadHighscoreWithFile(t *testing.T) {
	WriteHighscore(datafile, createTestScore())
	defer os.Remove(datafile)

	score, err := ReadHighscore(datafile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(score.Scores) != 1 {
		t.Errorf("Scores has more than one entry")
	}
}

func TestCheckScoreWithFile(t *testing.T) {
	WriteHighscore(datafile, createTestScore())
	defer os.Remove(datafile)

	score, err := ReadHighscore(datafile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if score.CheckScore(200) == -1 {
		t.Errorf("Scores check for new record has failed")
	}
	if score.CheckScore(50) != -1 {
		t.Errorf("Scores check has failed")
	}
}

func TestInsertScoreWithFile(t *testing.T) {
	WriteHighscore(datafile, createTestScore())
	defer os.Remove(datafile)

	score, err := ReadHighscore(datafile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	score.InsertScore(*NewScore(200, 0, "Test2"))
	if score.PrintScores() != "200,100," {
		t.Errorf("Insert new record failed")
	}

	score.InsertScore(*NewScore(50, 0, "Test3"))
	if score.PrintScores() != "200,100,50," {
		t.Errorf("Insert new score failed")
	}
}

func createTestScore() *HighScore {
	testscore := NewHighScore()
	testscore.Scores = append(testscore.Scores, *NewScore(100, 0, "Test"))
	return testscore
}
