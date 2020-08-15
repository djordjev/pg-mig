package utils

import (
	"encoding/json"

	"github.com/spf13/afero"
)

// Fs Setting up file system
var Fs = afero.NewOsFs()

// Config JSON type for storing database configuration
type Config struct {
	DbName      string `json:"db_name"`
	Path        string `json:"path"`
	DbURL       string `json:"db_url"`
	Credentials string `json:"credentials"`
	Port        int    `json:"port"`
	SSL         string `json:"ssl_mode"`
}

const configFileName = "pgmig.config.json"

// Store - saves configuration in json file
func (config *Config) Store() error {
	afs := &afero.Afero{Fs: Fs}
	exists, err := afs.Exists(configFileName)

	if err != nil {
		return err
	}

	if exists {
		err = afs.Remove(configFileName)

		if err != nil {
			return err
		}
	}

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = afs.WriteFile(configFileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Load - reads previously stored config file from current dir
func (config *Config) Load() error {
	afs := &afero.Afero{Fs: Fs}
	data, err := afs.ReadFile(configFileName)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}

	return nil
}
