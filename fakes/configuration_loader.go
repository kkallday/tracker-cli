package fakes

import "github.com/kkelani/tracker-cli/config"

type ConfigurationLoader struct {
	LoadCall struct {
		CallCount int
		Receives  struct {
			PathToConfig string
		}
		Returns struct {
			Configuration config.Configuration
			Error         error
		}
	}
}

func (c *ConfigurationLoader) Load(pathToConfig string) (config.Configuration, error) {
	c.LoadCall.CallCount++
	c.LoadCall.Receives.PathToConfig = pathToConfig
	return c.LoadCall.Returns.Configuration, c.LoadCall.Returns.Error
}
