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
	pathToConfigFile := writeConfigFile(t, `{"token": "some-token", "project_id": 54, "api_endpoint_override": "http://www.some-endpoint.com"}`)
	defer os.Remove(pathToConfigFile)

	configLoader := application.NewConfigurationLoader()
	configuration, err := configLoader.Load(path.Dir(pathToConfigFile))
	if err != nil {
		t.Errorf("Load() returned an unexpected error %q", err.Error())
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
		t.Errorf("Load() did not return a \"no such file or directory\" error was %q", err.Error())
	}
}

func TestLoadReturnsErrorWhenJSONDecodingFails(t *testing.T) {
	pathToConfigFile := writeConfigFile(t, `not-valid-json`)
	defer os.Remove(pathToConfigFile)

	configLoader := application.NewConfigurationLoader()
	_, err := configLoader.Load(path.Dir(pathToConfigFile))

	if !strings.Contains(err.Error(), "invalid character") {
		t.Errorf("Load() did not return a \"json decoding\" error %q", err.Error())
	}
}

func TestLoadReturnsErrorWhenFileIsMissingRequiredValues(t *testing.T) {
	pathToConfigFile := writeConfigFile(t, `{"api_endpoint_override": "http://www.some-endpoint.com"}`)
	defer os.Remove(pathToConfigFile)

	configLoader := application.NewConfigurationLoader()
	_, actualErr := configLoader.Load(path.Dir(pathToConfigFile))

	expectedErr := errors.New("Configuration must contain a token and a project ID")

	if actualErr.Error() != expectedErr.Error() {
		t.Errorf("Load() returned error %q, expected %q", actualErr.Error(), expectedErr.Error())
	}
}

func writeConfigFile(t *testing.T, configJSON string) string {
	configDirPath, err := ioutil.TempDir("", "")
	if err != nil {
		t.Errorf("Failed to create temp directory %q", err.Error())
	}

	configFile, err := os.Create(filepath.Join(configDirPath, "config.json"))
	defer configFile.Close()
	if err != nil {
		t.Errorf("Failed to create temp file %q", err.Error())
	}

	_, err = configFile.WriteString(configJSON)
	if err != nil {
		t.Errorf("Failed to write to temp file %q", err.Error())
	}

	return configFile.Name()
}
