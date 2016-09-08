package trackerapi_test

import (
	"github.com/kkallday/tracker-cli/trackerapi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client Provider", func() {
	Describe("Client", func() {
		It("returns client with specified api url and token", func() {
			clientProvider := trackerapi.NewClientProvider("http://www.some-tracker-api.com")
			actualClient := clientProvider.Client(12345, "some-token")

			expectedClient := trackerapi.TrackerClient{
				URL:       "http://www.some-tracker-api.com",
				Token:     "some-token",
				ProjectID: 12345,
			}

			Expect(actualClient).To(Equal(expectedClient))
		})
	})
})
