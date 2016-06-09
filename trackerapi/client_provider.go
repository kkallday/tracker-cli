package trackerapi

type ClientProvider struct {
}

func NewClientProvider() ClientProvider {
	return ClientProvider{}
}

func (ClientProvider) Client(token, url string) Client {
	if url == "" {
		url = "https://www.pivotaltracker.com"
	}

	return TrackerClient{
		URL:   url,
		Token: token,
	}
}
