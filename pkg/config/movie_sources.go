package config

import (
	"encoding/json"
	"os"
)

func LoadMovieSources(path string) (map[string]string, error) {
	file, err := os.Open(path)
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
