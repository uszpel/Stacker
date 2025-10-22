package main

import (
	"encoding/json"
	"log"
	"os"
)

type Score struct {
	Rank  int    `json:"rank"`
	Score int    `json:"score"`
	Lines int    `json:"lines"`
	Name  string `json:"name"`
}

type HighScore struct {
	Scores []Score `json:"scores"`
}

func readHighscore(datafile string) (*HighScore, error) {
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
		return data, nil
	} else {
		log.Printf("Highscore file %v not found", datafile)
	}

	return &HighScore{
		Scores: make([]Score, 0),
	}, nil
}

func writeHighscore(datafile string, highScore *HighScore) error {
	result, err := json.Marshal(highScore)
	if err != nil {
		return err
	}

	//TODO: encrypt data
	err = os.WriteFile(datafile, result, 0466)
	if err != nil {
		return err
	}
	return nil
}
