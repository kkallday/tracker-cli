package application

import (
	"fmt"
	"os"
	"strconv"
)

var getenv = os.Getenv

type Configuration struct {
	Token     string `json:"token"`
	ProjectID int    `json:"project_id"`
}

type ConfigurationLoader struct{}

func NewConfigurationLoader() ConfigurationLoader {
	return ConfigurationLoader{}
}

func (ConfigurationLoader) Load() (Configuration, error) {
	projectIDStr := getenv("PROJECT_ID")
	if projectIDStr == "" {
		return Configuration{}, fmt.Errorf("PROJECT_ID is required.")
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		return Configuration{}, fmt.Errorf("%q is not a valid PROJECT_ID. A number is required.", projectIDStr)
	}

	token := getenv("TOKEN")
	return Configuration{
		ProjectID: projectID,
		Token:     token,
	}, nil
}
