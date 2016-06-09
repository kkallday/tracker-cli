package fakes

import "github.com/kkelani/tracker-cli/application"

type ConfigurationLoader struct {
	LoadCall struct {
		CallCount int
		Receives  struct {
			PathToConfig string
		}
		Returns struct {
			Configuration application.Configuration
			Error         error
		}
	}
}

func (c *ConfigurationLoader) Load(pathToConfig string) (application.Configuration, error) {
	c.LoadCall.CallCount++
	c.LoadCall.Receives.PathToConfig = pathToConfig
	return c.LoadCall.Returns.Configuration, c.LoadCall.Returns.Error
}
