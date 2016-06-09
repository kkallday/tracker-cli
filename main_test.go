package main_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"
)

func TestMainShowsInFlightStories(t *testing.T) {
	buildTrackerCLI(t)

	server := startTestServer(t)
	defer server.Close()

	configFile := writeConfigFile(t, server.URL)
	defer os.Remove(configFile.Name())

	cmd := exec.Command("tracker-cli", "--config-dir", path.Dir(configFile.Name()))

	stdout := bytes.NewBuffer([]byte{})
	cmd.Stdout = stdout

	err := cmd.Run()
	if err != nil {
		t.Errorf("tracker-cli execution failed. Err: %v", err)
	}

	actualStdoutContentString := stdout.String()
	expectedStdoutContentString := loadFixture(t, "stories.stdout")
	if actualStdoutContentString != expectedStdoutContentString {
		fmt.Println(actualStdoutContentString)
		t.Errorf("Stdout is %q, did not match the fixture", actualStdoutContentString)
	}
}

func buildTrackerCLI(t *testing.T) {
	err := exec.Command("go", "install", "github.com/kkelani/tracker-cli").Run()
	if err != nil {
		t.Errorf("Failed to install tracker-cli\nError: %v", err)
	}
}

func startTestServer(t *testing.T) *httptest.Server {
	fixture := loadFixture(t, "project-stories.json")
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, fixture)
	}))

	return testServer
}

func writeConfigFile(t *testing.T, url string) *os.File {
	configDirPath, err := ioutil.TempDir("", "")
	if err != nil {
		t.Errorf("Failed to create temp directory %v", err)
	}

	configFile, err := os.Create(filepath.Join(configDirPath, "config.json"))
	if err != nil {
		t.Errorf("Failed to create temp file %v", err)
	}

	fileContents := fmt.Sprintf(`{"token": "some-token", "project_id": 105, "api_endpoint_override": %q}`, url)
	_, err = configFile.WriteString(fileContents)
	if err != nil {
		t.Errorf("Failed to write to temp file %v", err)
	}

	return configFile
}

func loadFixture(t *testing.T, fileName string) string {
	fixtureFileContents, err := ioutil.ReadFile(fmt.Sprintf("fixtures/%s", fileName))
	if err != nil {
		t.Errorf("Failed to load fixture file. %v", err)
	}

	return string(fixtureFileContents)
}
