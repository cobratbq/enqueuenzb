package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Url string
	Key string
	Delete bool
}

func readConfig(filePath string) (*Config, error) {
	path := os.ExpandEnv(filePath)
	// Read config file
	configFile, err := os.Open(os.ExpandEnv(path))
	if err != nil {
		return nil, fmt.Errorf("Cannot read config file '%s': %s\n", path, err.Error())
	}
	defer configFile.Close()
	// Decode json content
	config := new(Config)
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("Failed to decode json config file '%s': %s\n", path, err.Error())
	}
	return config, nil
}
