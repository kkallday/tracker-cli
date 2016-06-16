package application_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kkallday/tracker-cli/application"
	"github.com/kkallday/tracker-cli/fakes"
	"github.com/kkallday/tracker-cli/trackerapi"
)

var (
	fakeConfigurationLoader *fakes.ConfigurationLoader
	fakeClient              *fakes.Client
	fakeClientProvider      *fakes.ClientProvider
	fakeLogger              *fakes.Logger
)

func injectFakes() application.App {
	fakeConfigurationLoader = &fakes.ConfigurationLoader{}
	fakeClient = &fakes.Client{}
	fakeClientProvider = &fakes.ClientProvider{}
	fakeLogger = &fakes.Logger{}

	fakeClientProvider.ClientCall.Returns.Client = fakeClient

	return application.NewApp(fakeClientProvider, fakeConfigurationLoader, fakeLogger)
}

func TestRunRetrievesValuesFromConfigurationFile(t *testing.T) {
	app := injectFakes()

	err := app.Run("/dir/containing/config")

	if err != nil {
		t.Errorf("Run() returned an unexpected error %v", err)
	}

	actualCallCount := fakeConfigurationLoader.LoadCall.CallCount
	expectedCallCount := 1
	if actualCallCount != expectedCallCount {
		t.Errorf("Run() called configurationLoader.Load() %d time(s), expected %d time(s)", actualCallCount, expectedCallCount)
	}

	actualPathToConfig := fakeConfigurationLoader.LoadCall.Receives.PathToConfig
	expectedPathToConfig := "/dir/containing/config"
	if actualPathToConfig != expectedPathToConfig {
		t.Errorf("Run() called configurationLoader.Load(%q), expected configurationLoader.Load(%q)", actualPathToConfig, expectedPathToConfig)
	}
}

func TestRunReturnsErrorWhenConfigurationLoaderFails(t *testing.T) {
	app := injectFakes()
	fakeConfigurationLoader.LoadCall.Returns.Error = errors.New("load failed")

	actualErr := app.Run("")

	if actualErr == nil {
		t.Error("Run() did not return an expected error")
	}

	expectedErr := errors.New("load failed")
	if actualErr.Error() != expectedErr.Error() {
		t.Errorf("Run() returned error %q, expected error %q", actualErr.Error(), expectedErr.Error())
	}
}

func TestRunInitializesClientWithConfiguration(t *testing.T) {
	app := injectFakes()
	fakeConfigurationLoader.LoadCall.Returns.Configuration = application.Configuration{
		Token:               "some-token",
		APIEndpointOverride: "http://www.some-other-tracker.com",
	}

	err := app.Run("")

	if err != nil {
		t.Errorf("Run() returned an unexpected error %v", err)
	}

	actualCallCount := fakeClientProvider.ClientCall.CallCount
	expectedCallCount := 1
	if actualCallCount != expectedCallCount {
		t.Errorf("Run() called clientProvider.Client() %d time(s), expected %d time(s)", actualCallCount, expectedCallCount)
	}

	actualURL := fakeClientProvider.ClientCall.Receives.URL
	expectedURL := "http://www.some-other-tracker.com"
	actualToken := fakeClientProvider.ClientCall.Receives.Token
	expectedToken := "some-token"
	if actualURL != expectedURL || actualToken != expectedToken {
		t.Errorf("Run() called clientProvider.Client(%q, %q), expected clientProvider.Client(%q, %q)", actualURL, actualToken, expectedURL, expectedToken)
	}
}

func TestRunClientRetrievesProjectStories(t *testing.T) {
	app := injectFakes()
	fakeConfigurationLoader.LoadCall.Returns.Configuration = application.Configuration{
		ProjectID: 28,
	}

	err := app.Run("")

	if err != nil {
		t.Errorf("Run() returned an unexpected error %v", err)
	}

	actualCallCount := fakeClient.ProjectStoriesCall.CallCount
	expectedCallCount := 1
	if actualCallCount != expectedCallCount {
		t.Errorf("Run() called client.ProjectStories() %d time(s), expected %d time(s)", actualCallCount, expectedCallCount)
	}

	actualProjectID := fakeClient.ProjectStoriesCall.Receives.ProjectID
	expectedProjectID := 28
	if actualProjectID != expectedProjectID {
		t.Errorf("Run() called client.ProjectStories(%d), expected client.ProjectStories(%d)", actualProjectID, expectedProjectID)
	}
}

func TestRunClientReturnsErrorWhenRetrievingProjectStoriesFails(t *testing.T) {
	app := injectFakes()
	fakeClient.ProjectStoriesCall.Returns.Error = errors.New("failed to retrieve project stories")

	actualErr := app.Run("")
	if actualErr == nil {
		t.Error("Run() did not return an expected error")
	}

	expectedErr := errors.New("failed to retrieve project stories")
	if actualErr.Error() != expectedErr.Error() {
		t.Errorf("Run() returned error %q, expected error %q", actualErr.Error(), expectedErr.Error())
	}
}

func TestRunClientWritesTitleToLogger(t *testing.T) {
	app := injectFakes()

	err := app.Run("")

	if err != nil {
		t.Errorf("Run() returned an unexpected error %v", err)
	}

	actualCallCount := fakeLogger.LogCall.CallCount
	expectedCallCount := 1
	if actualCallCount != expectedCallCount {
		t.Errorf("Run() called logger.LogStories %d time(s), expected %d time(s)", actualCallCount, expectedCallCount)
	}

	actualLogMessage := fakeLogger.LogCall.Receives.Message
	expectedLogMessage := "Stories in-flight:"
	if actualLogMessage != expectedLogMessage {
		t.Errorf("Run() called logger.LogStories(%q), expected logger.LogStories(%q)", actualLogMessage, expectedLogMessage)
	}
}

func TestRunClientWritesStoriesToLogger(t *testing.T) {
	app := injectFakes()
	fakeClient.ProjectStoriesCall.Returns.Stories = []trackerapi.Story{
		{109832, "feature", "User can do X", 2},
		{201294, "bug", "something is wrong", 0},
		{838312, "chore", "this is a chore", 0},
	}

	err := app.Run("")

	if err != nil {
		t.Errorf("Run() returned an unexpected error %v", err)
	}

	actualCallCount := fakeLogger.LogStoriesCall.CallCount
	expectedCallCount := 1
	if actualCallCount != expectedCallCount {
		t.Errorf("Run() called logger.LogStories() %d time(s), expected %d time(s)", actualCallCount, expectedCallCount)
	}

	actualStories := fakeLogger.LogStoriesCall.Receives.Stories
	expectedStories := []trackerapi.Story{
		{109832, "feature", "User can do X", 2},
		{201294, "bug", "something is wrong", 0},
		{838312, "chore", "this is a chore", 0},
	}

	if !reflect.DeepEqual(actualStories, expectedStories) {
		t.Errorf("Run() called \nlogger.LogStories(%+v)\nexpected \nlogger.LogStories(%+v)", actualStories, expectedStories)
	}
}
