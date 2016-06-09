package application_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/kkelani/tracker-cli/application"
	"github.com/kkelani/tracker-cli/trackerapi"
)

func TestLoggerLogStoriesWritesStories(t *testing.T) {
	fakeWriter := bytes.NewBuffer([]byte{})

	logger := application.NewLogger(fakeWriter)

	logger.LogStories(
		trackerapi.Story{109832, "feature", "User can signup", 2},
		trackerapi.Story{201294, "bug", "something is wrong", 0},
	)

	actualContents, err := fakeWriter.ReadBytes(byte('\n'))
	if err != nil {
		t.Errorf("ReadBytes() returned unexpected error %v", err)
	}
	expectedContents := []byte("Feature   #109832   Pts: 2 User can signup\n")
	if !reflect.DeepEqual(actualContents, expectedContents) {
		t.Errorf("LogStories() wrote %q on line 1, expected %q\n",
			string(actualContents), string(expectedContents))
	}

	actualContents, err = fakeWriter.ReadBytes(byte('\n'))
	if err != nil {
		t.Errorf("ReadBytes() returned unexpected error %v", err)
	}
	expectedContents = []byte("Bug       #201294   something is wrong\n")
	if !reflect.DeepEqual(actualContents, expectedContents) {
		t.Errorf("LogStories() wrote %q on line 2, expected %q\n",
			string(actualContents), string(expectedContents))
	}
}

func TestLoggerLogWritesMessage(t *testing.T) {
	fakeWriter := bytes.NewBuffer([]byte{})

	logger := application.NewLogger(fakeWriter)
	logger.Log("some message")

	actualContents, err := fakeWriter.ReadBytes(byte('\n'))
	if err != nil {
		t.Errorf("ReadBytes() returned unexpected error %v", err)
	}
	expectedContents := []byte("some message\n")
	if !reflect.DeepEqual(actualContents, expectedContents) {
		t.Errorf("LogStories() wrote %q on line 1, expected %q\n",
			string(actualContents), string(expectedContents))
	}
}
