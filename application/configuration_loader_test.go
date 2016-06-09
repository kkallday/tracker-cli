package application_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kkelani/tracker-cli/application"
)

func TestLoadReturnsConfiguration(t *testing.T) {
	configFile := writeConfigFile(t, `{"token": "some-token", "project_id": 54, "api_endpoint_override": "http://www.some-endpoint.com"}`)
	defer os.Remove(configFile.Name())

	configLoader := application.NewConfigurationLoader()
	configuration, err := configLoader.Load(path.Dir(configFile.Name()))
	if err != nil {
		t.Errorf("Load() returned an unexpected error %v", err)
	}

	actualConfig := configuration
	expectedConfig := application.Configuration{
		Token:               "some-token",
		ProjectID:           54,
		APIEndpointOverride: "http://www.some-endpoint.com",
	}

	if actualConfig != expectedConfig {
		t.Errorf("Load() returned %+v, expected %+v", actualConfig, expectedConfig)
	}
}

func TestLoadReturnsErrorWhenFileOpeningFails(t *testing.T) {
	configLoader := application.NewConfigurationLoader()
	_, err := configLoader.Load("/non/existent/dir/path")

	if !strings.Contains(err.Error(), "/non/existent/dir/path/config.json: no such file or directory") {
		t.Errorf("Load() did not return a \"no such file or directory\" error was %v", err)
	}
}

func TestLoadReturnsErrorWhenJSONDecodingFails(t *testing.T) {
	configFile := writeConfigFile(t, `not-valid-json`)
	defer os.Remove(configFile.Name())

	configLoader := application.NewConfigurationLoader()
	_, err := configLoader.Load(path.Dir(configFile.Name()))

	if !strings.Contains(err.Error(), "invalid character") {
		t.Errorf("Load() did not return a \"json\" error %v", err)
	}
}

func TestLoadReturnsErrorWhenFileIsMissingRequiredValues(t *testing.T) {
	configFile := writeConfigFile(t, `{"api_endpoint_override": "http://www.some-endpoint.com"}`)
	defer os.Remove(configFile.Name())

	configLoader := application.NewConfigurationLoader()
	_, err := configLoader.Load(path.Dir(configFile.Name()))

	actualError := err
	expectedError := errors.New("Configuration must contain a token and a project ID")

	if actualError.Error() != expectedError.Error() {
		t.Errorf("Load() returned error %q, expected %q",
			actualError.Error(), expectedError.Error())
	}
}

func writeConfigFile(t *testing.T, configJSON string) *os.File {
	configDirPath, err := ioutil.TempDir("", "")
	if err != nil {
		t.Errorf("Failed to create temp directory %v", err)
	}

	configFile, err := os.Create(filepath.Join(configDirPath, "config.json"))
	if err != nil {
		t.Errorf("Failed to create temp file %v", err)
	}

	_, err = configFile.WriteString(configJSON)
	if err != nil {
		t.Errorf("Failed to write to temp file %v", err)
	}

	return configFile
}
