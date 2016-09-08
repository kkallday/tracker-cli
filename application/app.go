package application

import "github.com/kkallday/tracker-cli/trackerapi"

type App struct {
	clientProvider      clientProvider
	configurationLoader configurationLoader
	logger              logger
}

type clientProvider interface {
	Client(projectID int, token string) trackerapi.Client
}

type configurationLoader interface {
	Load() (Configuration, error)
}

type logger interface {
	LogStories(stories ...trackerapi.Story)
	Log(message string)
}

func NewApp(logger logger, clientProvider clientProvider, configurationLoader configurationLoader) App {
	return App{
		clientProvider:      clientProvider,
		configurationLoader: configurationLoader,
		logger:              logger,
	}
}

func (a App) Run() error {
	cfg, err := a.configurationLoader.Load()
	if err != nil {
		return err
	}

	client := a.clientProvider.Client(cfg.ProjectID, cfg.Token)
	stories, err := client.ProjectStories()
	if err != nil {
		return err
	}

	a.logger.Log("Stories in-flight:")
	a.logger.LogStories(stories...)

	return nil
}
