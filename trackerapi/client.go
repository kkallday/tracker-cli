package trackerapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Story struct {
	Id         int
	Story_type string
	Name       string
	Estimate   int
}

type Client interface {
	ProjectStories(projectId int) ([]Story, error)
}

type TrackerClient struct {
	URL   string
	Token string
}

func NewClient(url, token string) Client {
	return TrackerClient{
		URL:   url,
		Token: token,
	}
}

func (c TrackerClient) ProjectStories(projectId int) ([]Story, error) {
	targetURL := fmt.Sprintf("%s/services/v5/projects/%d/stories?with_state=started",
		c.URL, projectId)

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return []Story{}, err
	}
	req.Header.Set("X-TrackerToken", c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []Story{}, err
	}

	jsonResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		return []Story{}, fmt.Errorf("bad response: %s", jsonResp)
	}

	var stories []Story
	err = json.Unmarshal([]byte(jsonResp), &stories)
	if err != nil {
		return []Story{}, err
	}

	return stories, nil
}
