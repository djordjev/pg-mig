package filesystem

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/afero"
)

// Config JSON type for storing database configuration
type Config struct {
	DbName      string `json:"db_name"`
	Path        string `json:"path"`
	DbURL       string `json:"db_url"`
	Credentials string `json:"credentials"`
	Port        int    `json:"port"`
	SSL         string `json:"ssl_mode"`

	Filesystem afero.Fs `json:"-"`
}

const configFileName = "pgmig.config.json"

// Store - saves configuration in json file
func (config *Config) Store() error {
	afs := &afero.Afero{Fs: config.Filesystem}
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
	afs := &afero.Afero{Fs: config.Filesystem}
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

// GetConnectionString returns string for connecting on DB
func (config *Config) GetConnectionString() (string, error) {

	if config.Credentials == "" || config.DbName == "" || config.DbURL == "" || config.Port == 0 {
		return "", errors.New("invalid data in config file")
	}

	connectionString := fmt.Sprintf("postgres://%s@%s:%d/%s?sslmode=%s", config.Credentials, config.DbURL, config.Port, config.DbName, config.SSL)

	return connectionString, nil
}
