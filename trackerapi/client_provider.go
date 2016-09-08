package trackerapi

type ClientProvider struct {
	apiURL string
}

func NewClientProvider(apiURL string) ClientProvider {
	return ClientProvider{
		apiURL: apiURL,
	}
}

func (c ClientProvider) Client(projectID int, token string) Client {
	return TrackerClient{
		URL:       c.apiURL,
		Token:     token,
		ProjectID: projectID,
	}
}
