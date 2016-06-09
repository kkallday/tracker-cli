package trackerapi_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/kkelani/tracker-cli/trackerapi"
)

func TestClientProjectStoriesReturnsStories(t *testing.T) {
	testServer := startTestServer(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, storiesJSONFixture(t))
	})
	defer testServer.Close()

	client := trackerapi.NewClient(testServer.URL, "")
	actualStories, err := client.ProjectStories(1)

	if err != nil {
		t.Errorf("ProjectStories() returned error %v", err)
	}

	expectedStories := []trackerapi.Story{
		{Id: 1091909, Story_type: "feature", Name: "User can signup", Estimate: 1},
		{Id: 1909283, Story_type: "feature", Name: "User can create a todo list", Estimate: 2},
		{Id: 1032183, Story_type: "chore", Name: "Refactor app startup script", Estimate: 0},
		{Id: 2308423, Story_type: "bug", Name: "Signup success message should go away after navigating away", Estimate: 0},
	}

	if !reflect.DeepEqual(actualStories, expectedStories) {
		t.Errorf("ProjectStories() = %v, expected %v", actualStories, expectedStories)
	}
}

func TestClientProjectStoriesRequestContainsCorrectPathAndQueryAndHeader(t *testing.T) {
	testServer := startTestServer(func(w http.ResponseWriter, r *http.Request) {
		actualPath := r.URL.Path
		expectedPath := "/services/v5/projects/6/stories"
		if actualPath != expectedPath {
			t.Errorf("Path was %q, expected %q", actualPath, expectedPath)
		}

		actualQuery := r.URL.Query().Get("with_state")
		expectedQuery := "started"
		if actualQuery != expectedQuery {
			t.Errorf("GET query value for \"with_state\" was %q, expected %q", actualQuery, expectedQuery)
		}

		actualTrackerToken := r.Header.Get("X-TrackerToken")
		expectedTrackerToken := "some-tracker-api-token"

		if actualTrackerToken != expectedTrackerToken {
			t.Errorf("Request header X-TrackerToken value was %q, expected %q",
				actualTrackerToken, expectedTrackerToken)
		}

		fmt.Fprint(w, "[]")
	})
	defer testServer.Close()

	client := trackerapi.NewClient(testServer.URL, "some-tracker-api-token")
	client.ProjectStories(6)
}

func TestClientReturnsErrorWhenRequestCreationFails(t *testing.T) {
	client := trackerapi.NewClient("http://%%%%%", "")

	_, err := client.ProjectStories(9)
	if err == nil {
		t.Error("ProjectStories() did not return an error, expected an error")
	}

	actualErr := err
	expectedErr := errors.New("parse http://%%%%%/services/v5/projects/9/stories?with_state=started: invalid URL escape \"%%%\"")
	if actualErr.Error() != expectedErr.Error() {
		t.Errorf("ProjectStories() returned error %q, expected %q", actualErr.Error(), expectedErr.Error())
	}
}

func TestClientReturnsErrorWhenRequestFails(t *testing.T) {
	client := trackerapi.NewClient("UNKNOWN-PROTOCOL://foo.com", "")

	_, err := client.ProjectStories(4)
	if err == nil {
		t.Error("ProjectStories() did not return an error, expected an error")
	}

	actualErr := err
	expectedErr := errors.New("Get unknown-protocol://foo.com/services/v5/projects/4/stories?with_state=started: unsupported protocol scheme \"unknown-protocol\"")
	if actualErr.Error() != expectedErr.Error() {
		t.Errorf("ProjectStories() returned error %q, expected %q", actualErr.Error(), expectedErr.Error())
	}
}

func TestClientReturnsErrorWhenJSONDecodingFails(t *testing.T) {
	testServer := startTestServer(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "not-valid-json")
	})
	defer testServer.Close()

	client := trackerapi.NewClient(testServer.URL, "")
	_, err := client.ProjectStories(1)

	if err == nil {
		t.Error("ProjectStories() did not return an error, expected an error")
	}

	actualErr := err
	expectedErr := errors.New("invalid character 'o' in literal null (expecting 'u')")
	if actualErr.Error() != expectedErr.Error() {
		t.Errorf("ProjectStories() returned error %q, expected %q", actualErr.Error(), expectedErr.Error())
	}
}

func TestClientReturnsErrorWhenRequestIsNotOK(t *testing.T) {
	testServer := startTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "an error occurred"}`))
	})
	defer testServer.Close()

	client := trackerapi.NewClient(testServer.URL, "")
	_, err := client.ProjectStories(1)

	expectedErr := errors.New(`bad response: {"message": "an error occurred"}`)
	if err.Error() != expectedErr.Error() {
		t.Errorf("ProjectStories() returned error %q, expected %q", err.Error(), expectedErr.Error())
	}
}

func startTestServer(handler func(http.ResponseWriter, *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

func storiesJSONFixture(t *testing.T) string {
	fixtureFileContents, err := ioutil.ReadFile("fixtures/stories.json")
	if err != nil {
		t.Errorf("Failed to load fixture file. %v", err)
	}

	return string(fixtureFileContents)
}
