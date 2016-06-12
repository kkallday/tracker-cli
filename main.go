package main

import (
	"os"

	"github.com/kkallday/tracker-cli/application"
	"github.com/kkallday/tracker-cli/trackerapi"
)

func main() {
	logger := application.NewLogger(os.Stdout)
	configurationLoader := application.NewConfigurationLoader()
	clientProvider := trackerapi.NewClientProvider()

	flagParser := application.NewFlagParser()
	cmdLineConfig, err := flagParser.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}

	app := application.NewApp(clientProvider, configurationLoader, logger)
	err = app.Run(cmdLineConfig.ConfigDir)
	if err != nil {
		panic(err)
	}

}
