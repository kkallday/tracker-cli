package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kkallday/tracker-cli/application"
	"github.com/kkallday/tracker-cli/trackerapi"
)

const (
	trackerAPIURL = "https://www.pivotaltracker.com"
)

func main() {
	logger := application.NewLogger(os.Stdout)
	configurationLoader := application.NewConfigurationLoader()

	apiURL := flag.String("api", trackerAPIURL, "url of api (defaults to Pivotal Tracker api)")
	flag.Parse()
	clientProvider := trackerapi.NewClientProvider(*apiURL)

	app := application.NewApp(logger, clientProvider, configurationLoader)

	err := app.Run()
	if err != nil {
		failWithError(err)
	}
}

func failWithError(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
