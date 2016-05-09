package trackerapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Story struct {
	Id         int
	Story_type string
	Name       string
	Estimate   int
	Owner_ids  []int
}

type Client struct {
	URL        string
	Token      string
	HttpClient http.Client
}

func (c Client) ProjectStories(projectId int) ([]Story, error) {
	targetURL := fmt.Sprintf("%s/services/v5/projects/%d/stories?with_state=started", c.URL, projectId)
	resp, err := c.HttpClient.Get(targetURL)
	if err != nil {
		return []Story{}, err
	}

	var stories []Story
	err = json.NewDecoder(resp.Body).Decode(&stories)
	if err != nil {
		return []Story{}, err
	}

	return stories, nil
}
