package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Score struct {
	Score int    `json:"score"`
	Lines int    `json:"lines"`
	Name  string `json:"name"`
	Date  string `json:"date"`
	IsNew bool   `json:"-"`
}

type HighScore struct {
	Scores []Score `json:"scores"`
}

func NewHighScore() *HighScore {
	return &HighScore{
		Scores: make([]Score, 0),
	}
}

func NewScore(score int, lines int, name string) *Score {
	return &Score{
		Score: score,
		Name:  name,
		Lines: lines,
		Date:  time.Now().UTC().String(),
	}
}

func (h HighScore) CheckScore(score int) int {
	result := -1
	for index, curScore := range h.Scores {
		if curScore.Score <= score {
			result = index
			break
		}
	}
	return result
}

func (h *HighScore) InsertScore(score Score) {
	index := h.CheckScore(score.Score)
	score.IsNew = true
	if index == -1 {
		h.Scores = append(h.Scores, score)
	} else {
		newScores := make([]Score, 0)
		newScores = append(newScores, h.Scores[:index]...)
		newScores = append(newScores, score)
		newScores = append(newScores, h.Scores[index:]...)
		h.Scores = newScores
	}
	if len(h.Scores) > 10 {
		h.Scores = h.Scores[:10]
	}
}

func (h HighScore) PrintScores() string {
	result := ""
	for _, curScore := range h.Scores {
		result += fmt.Sprintf("%d,", curScore.Score)
	}
	return result
}

func ReadHighscore(datafile string) (*HighScore, error) {
	if _, err := os.Stat(datafile); err == nil {
		result, err := os.ReadFile(datafile)
		if err != nil {
			return nil, err
		}

		//TODO: Decrypt data
		data := &HighScore{}
		err = json.Unmarshal(result, &data)
		if err != nil {
			return nil, err
		}
		//TODO: Check order
		return data, nil
	} else {
		log.Printf("Highscore file %v not found", datafile)
	}

	return NewHighScore(), nil
}

func WriteHighscore(datafile string, highScore *HighScore) error {
	result, err := json.Marshal(highScore)
	if err != nil {
		return err
	}

	os.Remove(datafile)
	//TODO: encrypt data
	err = os.WriteFile(datafile, result, 0666)
	if err != nil {
		return err
	}
	return nil
}
