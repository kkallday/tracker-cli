package application_test

import (
	"errors"
	"testing"

	"github.com/kkallday/tracker-cli/application"
)

func TestParseReturnsCommandLineArgs(t *testing.T) {
	flagParser := application.NewFlagParser()
	actualCommandLineArgs, err := flagParser.Parse([]string{"--config-dir", "/path/to/config"})

	if err != nil {
		t.Errorf("Parse() return an unexpected error: %q", err.Error())
	}

	expectedCommandLineArgs := application.CommandLineArgs{ConfigDir: "/path/to/config"}
	if actualCommandLineArgs != expectedCommandLineArgs {
		t.Errorf("Parse() returned %+v, expected %+v", actualCommandLineArgs, expectedCommandLineArgs)
	}
}

func TestParseReturnsErrorWhenConfigDirIsMissing(t *testing.T) {
	flagParser := application.NewFlagParser()
	_, actualErr := flagParser.Parse([]string{"--some-arg", "some-value"})

	if actualErr == nil {
		t.Error("Parse() did not return an expected error")
	}

	expectedErr := errors.New("missing required flag --config-dir")
	if actualErr.Error() != expectedErr.Error() {
		t.Errorf("Parse() returned error %q, expected error %q", actualErr.Error(), expectedErr.Error())
	}
}
