package trackerapi_test

import (
	"github.com/kkallday/tracker-cli/trackerapi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client Provider", func() {
	Describe("Client", func() {
		It("returns client", func() {
			clientProvider := trackerapi.NewClientProvider()

			actualClient := clientProvider.Client("some-token", "http://www.some-tracker-api.com")

			expectedClient := trackerapi.TrackerClient{
				URL:   "http://www.some-tracker-api.com",
				Token: "some-token",
			}

			Expect(actualClient).To(Equal(expectedClient))
		})

		Context("when api endpoint is not specified", func() {
			It("returns a client with default URL", func() {
				clientProvider := trackerapi.NewClientProvider()

				actualClient := clientProvider.Client("some-token", "")

				expectedClient := trackerapi.TrackerClient{
					URL:   "https://www.pivotaltracker.com",
					Token: "some-token",
				}

				Expect(actualClient).To(Equal(expectedClient))
			})
		})
	})
})
