package fakes

import "github.com/kkallday/tracker-cli/trackerapi"

type Client struct {
	ProjectStoriesCall struct {
		CallCount int
		Receives  struct {
			ProjectID int
		}
		Returns struct {
			Stories []trackerapi.Story
			Error   error
		}
	}
}

func (c *Client) ProjectStories(projectID int) ([]trackerapi.Story, error) {
	c.ProjectStoriesCall.CallCount++
	c.ProjectStoriesCall.Receives.ProjectID = projectID
	return c.ProjectStoriesCall.Returns.Stories, c.ProjectStoriesCall.Returns.Error
}
