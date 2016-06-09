package application_test

import (
	"errors"
	"testing"

	"github.com/kkelani/tracker-cli/application"
)

func TestFlagParserParseReturnsCommandLineConfig(t *testing.T) {
	flagParser := application.NewFlagParser()
	actualCommandLineConfig, err := flagParser.Parse([]string{"--config-dir", "/path/to/config"})

	if err != nil {
		t.Errorf("Parse() return an unexpected error: %v", err)
	}

	expectedCommandLineConfig := application.CommandLineConfig{ConfigDir: "/path/to/config"}
	if actualCommandLineConfig != expectedCommandLineConfig {
		t.Errorf("Parse() returned %+v, expected %+v", actualCommandLineConfig, expectedCommandLineConfig)
	}
}

func TestFlagParserPaserReturnsErrorWhenConfigDirIsMissing(t *testing.T) {
	flagParser := application.NewFlagParser()
	_, actualErr := flagParser.Parse([]string{"--some-arg", "some-value"})

	if actualErr == nil {
		t.Errorf("Parse() did not return an expected error: %v", actualErr)
	}

	expectedErr := errors.New("missing required flag --config-dir")
	if actualErr.Error() != expectedErr.Error() {
		t.Errorf("Parse() returned error %v, expected error %v", actualErr, expectedErr)
	}
}

func TestFlagParserPaserReturnsErrorWhenConfigDirValueIsMissing(t *testing.T) {
	flagParser := application.NewFlagParser()
	_, actualErr := flagParser.Parse([]string{"--config-dir"})

	if actualErr == nil {
		t.Errorf("Parse() did not return an expected error: %v", actualErr)
	}

	expectedErr := errors.New("missing required flag --config-dir")
	if actualErr.Error() != expectedErr.Error() {
		t.Errorf("Parse() returned error %v, expected error %v", actualErr, expectedErr)
	}
}
