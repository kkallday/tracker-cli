package application_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kkelani/tracker-cli/application"
	"github.com/kkelani/tracker-cli/config"
	"github.com/kkelani/tracker-cli/fakes"
	"github.com/kkelani/tracker-cli/trackerapi"
)

var (
	fakeConfigurationLoader *fakes.ConfigurationLoader
	fakeClient              *fakes.Client
	fakeClientProvider      *fakes.ClientProvider
	fakeLogger              *fakes.Logger
)

func setupFakes() {
	fakeConfigurationLoader = &fakes.ConfigurationLoader{}
	fakeClient = &fakes.Client{}
	fakeClientProvider = &fakes.ClientProvider{}
	fakeLogger = &fakes.Logger{}

	fakeClientProvider.ClientCall.Returns.Client = fakeClient
}

func TestAppRunRetrievesValuesFromConfigurationFile(t *testing.T) {
	setupFakes()
	app := application.NewApp(fakeClientProvider, fakeConfigurationLoader, fakeLogger)
	err := app.Run("dir/containing/config")

	if err != nil {
		t.Errorf("Run() returned an unexpected error %v", err)
	}

	actualCallCount := fakeConfigurationLoader.LoadCall.CallCount
	expectedCallCount := 1
	if actualCallCount != expectedCallCount {
		t.Errorf("Run() called configurationLoader.Load() %d times, expected %d times", actualCallCount, expectedCallCount)
	}

	actualPathToConfig := fakeConfigurationLoader.LoadCall.Receives.PathToConfig
	expectedPathToConfig := "dir/containing/config"
	if actualPathToConfig != expectedPathToConfig {
		t.Errorf("Run() called configurationLoader.Load(%q), expected configurationLoader.Load(%q)", actualPathToConfig, expectedPathToConfig)
	}
}

func TestAppRunReturnsErrorWhenConfigurationLoaderFails(t *testing.T) {
	setupFakes()
	fakeConfigurationLoader.LoadCall.Returns.Error = errors.New("load failed")

	app := application.NewApp(fakeClientProvider, fakeConfigurationLoader, fakeLogger)
	err := app.Run("dir/containing/config")

	actualError := err
	expectedError := errors.New("load failed")
	if actualError.Error() != expectedError.Error() {
		t.Errorf("Run() returned error %q, expected error %q",
			actualError.Error(), expectedError.Error())
	}
}

func TestAppRunInitializesClientWithConfiguration(t *testing.T) {
	setupFakes()
	fakeConfigurationLoader.LoadCall.Returns.Configuration = config.Configuration{
		Token:               "some-token",
		APIEndpointOverride: "http://www.some-other-tracker.com",
	}

	app := application.NewApp(fakeClientProvider, fakeConfigurationLoader, fakeLogger)
	err := app.Run("dir/containing/config")

	if err != nil {
		t.Errorf("Run() returned an unexpected error %v", err)
	}

	actualCallCount := fakeClientProvider.ClientCall.CallCount
	expectedCallCount := 1
	if actualCallCount != expectedCallCount {
		t.Errorf("Run() called clientProvider.Client() %d times, expected %d time", actualCallCount, expectedCallCount)
	}

	actualURL := fakeClientProvider.ClientCall.Receives.URL
	expectedURL := "http://www.some-other-tracker.com"
	actualToken := fakeClientProvider.ClientCall.Receives.Token
	expectedToken := "some-token"
	if actualURL != expectedURL || actualToken != expectedToken {
		t.Errorf("Run() called clientProvider.Client(%q, %q), expected clientProvider.Client(%q, %q)",
			actualURL, actualToken, expectedURL, expectedToken)
	}
}

func TestAppRunClientRetrievesProjectStories(t *testing.T) {
	setupFakes()
	fakeConfigurationLoader.LoadCall.Returns.Configuration = config.Configuration{
		ProjectID: 28,
	}

	app := application.NewApp(fakeClientProvider, fakeConfigurationLoader, fakeLogger)
	err := app.Run("dir/containing/config")

	if err != nil {
		t.Errorf("Run() returned an unexpected error %v", err)
	}

	actualCallCount := fakeClient.ProjectStoriesCall.CallCount
	expectedCallCount := 1
	if actualCallCount != expectedCallCount {
		t.Errorf("Run() called client.ProjectStories() %d times, expected %d time", actualCallCount, expectedCallCount)
	}

	actualProjectID := fakeClient.ProjectStoriesCall.Receives.ProjectID
	expectedProjectID := 28
	if actualProjectID != expectedProjectID {
		t.Errorf("Run() called client.ProjectStories(%d), expected client.ProjectStories(%d)", actualProjectID, expectedProjectID)
	}
}

func TestAppRunClientReturnsErrorWhenRetrievingProjectStoriesFails(t *testing.T) {
	setupFakes()
	fakeClient.ProjectStoriesCall.Returns.Error = errors.New("failed to retrieve project stories")

	app := application.NewApp(fakeClientProvider, fakeConfigurationLoader, fakeLogger)
	err := app.Run("dir/containing/config")

	actualError := err
	expectedError := errors.New("failed to retrieve project stories")
	if actualError.Error() != expectedError.Error() {
		t.Errorf("Run() returned error %q, expected error %q", actualError.Error(), expectedError.Error())
	}
}

func TestAppRunClientWritesTitleToLogger(t *testing.T) {
	setupFakes()

	app := application.NewApp(fakeClientProvider, fakeConfigurationLoader, fakeLogger)
	err := app.Run("dir/containing/config")

	if err != nil {
		t.Errorf("Run() returned an unexpected error %v", err)
	}

	actualCallCount := fakeLogger.LogCall.CallCount
	expectedCallCount := 1
	if actualCallCount != expectedCallCount {
		t.Errorf("Run() called logger.LogStories %d time(s), expected %d time", actualCallCount, expectedCallCount)
	}

	actualLogMessage := fakeLogger.LogCall.Receives.Message
	expectedLogMessage := "Stories in-flight:"
	if actualLogMessage != expectedLogMessage {
		t.Errorf("Run() called logger.LogStories(%q), expected logger.LogStories(%q)", actualLogMessage, expectedLogMessage)
	}
}

func TestAppRunClientWritesStoriesToLogger(t *testing.T) {
	setupFakes()
	fakeClient.ProjectStoriesCall.Returns.Stories = []trackerapi.Story{
		{109832, "feature", "User can do X", 2},
		{201294, "bug", "something is wrong", 0},
		{838312, "chore", "this is a chore", 0},
	}

	app := application.NewApp(fakeClientProvider, fakeConfigurationLoader, fakeLogger)
	err := app.Run("dir/containing/config")

	if err != nil {
		t.Errorf("Run() returned an unexpected error %v", err)
	}

	actualCallCount := fakeLogger.LogStoriesCall.CallCount
	expectedCallCount := 1
	if actualCallCount != expectedCallCount {
		t.Errorf("Run() called logger.LogStories() %d times, expected %d time", actualCallCount, expectedCallCount)
	}

	actualStories := fakeLogger.LogStoriesCall.Receives.Stories
	expectedStories := []trackerapi.Story{
		{109832, "feature", "User can do X", 2},
		{201294, "bug", "something is wrong", 0},
		{838312, "chore", "this is a chore", 0},
	}

	if !reflect.DeepEqual(actualStories, expectedStories) {
		t.Errorf("Run() called logger.LogStories(%+v)\nexpected logger.LogStories(%+v)", actualStories, expectedStories)
	}
}

func TestAppRunReturnsErrorWhenLoggerFails(t *testing.T) {
	setupFakes()
	fakeLogger.LogStoriesCall.Returns.Error = errors.New("logging failed")

	app := application.NewApp(fakeClientProvider, fakeConfigurationLoader, fakeLogger)
	err := app.Run("dir/containing/config")

	actualError := err
	expectedError := errors.New("logging failed")
	if actualError.Error() != expectedError.Error() {
		t.Errorf("Run() returned error %q, expected error %q", actualError.Error(), expectedError.Error())
	}
}
