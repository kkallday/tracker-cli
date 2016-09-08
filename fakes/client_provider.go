package fakes

import "github.com/kkallday/tracker-cli/trackerapi"

type ClientProvider struct {
	ClientCall struct {
		CallCount int
		Receives  struct {
			ProjectID int
			Token     string
		}
		Returns struct {
			Client trackerapi.Client
		}
	}
}

func (c *ClientProvider) Client(projectID int, token string) trackerapi.Client {
	c.ClientCall.CallCount++
	c.ClientCall.Receives.ProjectID = projectID
	c.ClientCall.Receives.Token = token
	return c.ClientCall.Returns.Client
}
