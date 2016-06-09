package application

import (
	"encoding/json"
	"errors"
	"os"
	"path"
)

type Configuration struct {
	Token               string `json:"token"`
	ProjectID           int    `json:"project_id"`
	APIEndpointOverride string `json:"api_endpoint_override"`
}

type ConfigurationLoader struct{}

func NewConfigurationLoader() ConfigurationLoader {
	return ConfigurationLoader{}
}

func (ConfigurationLoader) Load(pathToConfigDir string) (Configuration, error) {
	pathToConfigFile := path.Join(pathToConfigDir, "config.json")
	configFile, err := os.OpenFile(pathToConfigFile, 0, os.FileMode(0644))
	if err != nil {
		return Configuration{}, err
	}

	var cfg Configuration
	err = json.NewDecoder(configFile).Decode(&cfg)
	if err != nil {
		return Configuration{}, err
	}

	if cfg.Token == "" || cfg.ProjectID == 0 {
		return Configuration{}, errors.New("Configuration must contain a token and a project ID")
	}

	return cfg, nil
}
