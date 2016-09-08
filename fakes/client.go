package fakes

import "github.com/kkallday/tracker-cli/trackerapi"

type Client struct {
	ProjectStoriesCall struct {
		CallCount int
		Returns   struct {
			Stories []trackerapi.Story
			Error   error
		}
	}
}

func (c *Client) ProjectStories() ([]trackerapi.Story, error) {
	c.ProjectStoriesCall.CallCount++
	return c.ProjectStoriesCall.Returns.Stories, c.ProjectStoriesCall.Returns.Error
}
