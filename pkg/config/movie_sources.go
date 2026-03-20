package config

import (
	"encoding/json"
	"os"
)

func LoadMovieSources() (map[string]string, error) {
	file, err := os.Open(Envs.MovieSourcesPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sources map[string]string
	if err := json.NewDecoder(file).Decode(&sources); err != nil {
		return nil, err
	}

	return sources, nil
}
