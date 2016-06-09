package trackerapi_test

import (
	"testing"

	"github.com/kkelani/tracker-cli/trackerapi"
)

func TestClientProviderReturnsClient(t *testing.T) {
	clientProvider := trackerapi.NewClientProvider()

	actualClient := clientProvider.Client("some-token", "http://www.some-tracker-api.com")

	expectedClient := trackerapi.TrackerClient{
		URL:   "http://www.some-tracker-api.com",
		Token: "some-token",
	}

	if actualClient != expectedClient {
		t.Errorf("Client() returned %+v, expected %+v", actualClient, expectedClient)
	}
}

func TestClientProviderReturnsClientWithDefaultEndpointOverride(t *testing.T) {
	clientProvider := trackerapi.NewClientProvider()

	actualClient := clientProvider.Client("some-token", "")

	expectedClient := trackerapi.TrackerClient{
		URL:   "https://www.pivotaltracker.com",
		Token: "some-token",
	}

	if actualClient != expectedClient {
		t.Errorf("Client() returned %+v, expected %+v", actualClient, expectedClient)
	}
}
