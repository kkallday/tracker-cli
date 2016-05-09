package trackerapi_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/kkelani/tracker-cli/trackerapi"
)

func TestClientProjectStoriesReturnsStories(t *testing.T) {
	httpClient, testServer := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(storiesFixture(t)))
	})
	defer testServer.Close()

	client := trackerapi.Client{
		URL:        testServer.URL,
		Token:      "some-tracker-api-token",
		HttpClient: httpClient,
	}

	actualStories, err := client.ProjectStories(101)

	if err != nil {
		t.Errorf("ProjectStories(101) returned error: %v", err)
	}

	expectedStories := []trackerapi.Story{
		{Id: 1091909, Story_type: "feature", Name: "Feature 1", Estimate: 1, Owner_ids: []int{5, 9}},
		{Id: 1909283, Story_type: "feature", Name: "Feature 2", Estimate: 2, Owner_ids: []int{5}},
		{Id: 1032183, Story_type: "chore", Name: "Chore 1", Estimate: 0, Owner_ids: []int{1, 7}},
		{Id: 2308423, Story_type: "bug", Name: "Bug 1", Estimate: 0, Owner_ids: []int{4}},
	}

	if !reflect.DeepEqual(actualStories, expectedStories) {
		t.Errorf("ProjectStories() = %v, expected %v", actualStories, expectedStories)
	}
}

func TestClientProjectStoriesMakesRequestToCorrectEndpointWithCorrectParams(t *testing.T) {
	var testServer *httptest.Server
	httpClient, testServer := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		requestURL := fmt.Sprintf("%s://%s", r.URL.Scheme, r.URL.Host)
		if requestURL != testServer.URL {
			t.Errorf("Server host is %q, expected %q", requestURL, testServer.URL)
		}

		if r.Method != "GET" {
			t.Errorf("Server method is %q, expected \"GET\"", r.Method)
		}

		expectedPath := "/services/v5/projects/101/stories"
		if r.URL.Path != expectedPath {
			t.Errorf("Server path is %q, expected %q", r.URL.Path, expectedPath)
		}

		expectedParams := "with_state=started"
		if r.URL.RawQuery != expectedParams {
			t.Error("Server params are %q, expected %q", "", expectedParams)
		}
	})
	defer testServer.Close()

	client := trackerapi.Client{
		URL:        testServer.URL,
		Token:      "some-tracker-api-token",
		HttpClient: httpClient,
	}

	client.ProjectStories(101)
}

func TestClientProjectStoriesReturnsErrorWhenHttpClientGetFails(t *testing.T) {
	httpClient, testServer := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {})
	defer testServer.Close()

	client := trackerapi.Client{
		URL:        "UNKNOWN-PROTOCOL://www.pivotaltracker.com",
		Token:      "some-tracker-api-token",
		HttpClient: httpClient,
	}
	_, err := client.ProjectStories(104)

	if err == nil {
		t.Error("ProjectStories() did not return expected error")
	}
}

func TestClientProjectStoriesReturnsErrorWhenJSONDecodeFails(t *testing.T) {
	httpClient, testServer := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "this-is-unparseable-json")
	})
	defer testServer.Close()

	client := trackerapi.Client{
		URL:        "http://www.pivotaltracker.com/project/stories",
		Token:      "some-tracker-api-token",
		HttpClient: httpClient,
	}
	_, err := client.ProjectStories(104)

	if err == nil {
		t.Error("ProjectStories() did not return expected error")
	}
}

func startTestServer(t *testing.T, handler func(http.ResponseWriter, *http.Request)) (*http.Client, *httptest.Server) {
	testServer := httptest.NewServer(http.HandlerFunc(handler))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(testServer.URL)
		},
	}
	httpClient := &http.Client{Transport: transport}

	return httpClient, testServer
}

func storiesFixture(t *testing.T) []byte {
	fixtureFileContents, err := ioutil.ReadFile("fixtures/2-features-1-chore-1-bug.json")
	if err != nil {
		t.Errorf("Failed to load fixture file. %v", err)
	}

	return fixtureFileContents
}
