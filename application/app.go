package application

import (
	"github.com/kkelani/tracker-cli/config"
	"github.com/kkelani/tracker-cli/trackerapi"
)

type App struct {
	clientProvider      clientProvider
	configurationLoader configurationLoader
	logger              logger
}

type clientProvider interface {
	Client(url, token string) trackerapi.Client
}

type configurationLoader interface {
	Load(pathToConfig string) (config.Configuration, error)
}

type logger interface {
	LogStories(stories ...trackerapi.Story) error
	Log(message string)
}

func NewApp(clientProvider clientProvider, configurationLoader configurationLoader, logger logger) App {
	return App{
		clientProvider:      clientProvider,
		configurationLoader: configurationLoader,
		logger:              logger,
	}
}

func (a App) Run(pathToConfig string) error {
	cfg, err := a.configurationLoader.Load(pathToConfig)
	if err != nil {
		return err
	}

	client := a.clientProvider.Client(cfg.Token, cfg.APIEndpointOverride)
	stories, err := client.ProjectStories(cfg.ProjectID)
	if err != nil {
		return err
	}

	a.logger.Log("Stories in-flight:")

	err = a.logger.LogStories(stories...)
	if err != nil {
		return err
	}

	return nil
}
