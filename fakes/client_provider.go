package fakes

import "github.com/kkallday/tracker-cli/trackerapi"

type ClientProvider struct {
	ClientCall struct {
		CallCount int
		Receives  struct {
			URL   string
			Token string
		}
		Returns struct {
			Client trackerapi.Client
		}
	}
}

func (c *ClientProvider) Client(token, url string) trackerapi.Client {
	c.ClientCall.CallCount++
	c.ClientCall.Receives.URL = url
	c.ClientCall.Receives.Token = token
	return c.ClientCall.Returns.Client
}
