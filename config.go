package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	StartAsMuted bool    `json:"start_as_muted"`
	SoundVolume  float64 `json:"sound_volume"`
}

func ReadConfig(file string) (*Config, error) {
	if _, err := os.Stat(file); err == nil {
		result, err := os.ReadFile(configfile)
		if err != nil {
			return nil, err
		}

		config := &Config{}
		err = json.Unmarshal(result, config)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
	return nil, fmt.Errorf("config file not found")
}
