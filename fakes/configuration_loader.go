package fakes

import "github.com/kkallday/tracker-cli/application"

type ConfigurationLoader struct {
	LoadCall struct {
		CallCount int
		Returns   struct {
			Configuration application.Configuration
			Error         error
		}
	}
}

func (c *ConfigurationLoader) Load() (application.Configuration, error) {
	c.LoadCall.CallCount++
	return c.LoadCall.Returns.Configuration, c.LoadCall.Returns.Error
}
