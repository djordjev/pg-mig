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
	NoColor     bool   `json:"no_color"`
}

const configFileName = "pgmig.config.json"

// StoreConfig - saves configuration in json file
func (fs *ImplFilesystem) StoreConfig(config Config) error {
	afs := &afero.Afero{Fs: fs.Fs}

	exists, err := afs.Exists(configFileName)

	if err != nil {
		return fmt.Errorf("filesystem error: unable to check if confg already exists %w", err)
	}

	if exists {
		err = afs.Remove(configFileName)

		if err != nil {
			return fmt.Errorf("filesystem error: unable to overwrite existing config %w", err)
		}
	}

	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("filesystem error: unable to serialize config data %w", err)
	}

	err = afs.WriteFile(configFileName, data, 0644)
	if err != nil {
		return fmt.Errorf("filesystem error: unable to write config file %w", err)
	}

	return nil
}

// LoadConfig - reads previously stored config file from current dir
func (fs *ImplFilesystem) LoadConfig() (Config, error) {
	afs := &afero.Afero{Fs: fs.Fs}
	config := Config{}

	data, err := afs.ReadFile(configFileName)

	if err != nil {
		return config, fmt.Errorf("filesystem error: unable to read config file %w", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("filesystem error: unable to desirialize existing configuration %s %w", string(data), err)
	}

	return config, nil
}

// GetConnectionString returns string for connecting on DB
func (config *Config) GetConnectionString() (string, error) {

	if config.Credentials == "" || config.DbName == "" || config.Port == 0 {
		return "", errors.New("filesystem error: invalid data in config file")
	}

	connectionString := fmt.Sprintf("postgres://%s@%s:%d/%s?sslmode=%s", config.Credentials, config.DbURL, config.Port, config.DbName, config.SSL)

	return connectionString, nil
}
