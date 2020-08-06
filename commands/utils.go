package commands

import (
	"encoding/json"
	"io/ioutil"
)

// Config JSON type for storing database configuration
type Config struct {
	DbName      string `json:"db_name"`
	Path        string `json:"path"`
	DbURL       string `json:"db_url"`
	Credentials string `json:"credentials"`
}

const configFileName = "pgmig.config.json"

// Store - saves configuration in json file
func (config *Config) Store() error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configFileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Load - reads previously stored config file from current dir
func (config *Config) Load() error {
	data, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}

	return nil
}
