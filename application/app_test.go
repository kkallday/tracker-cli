package application_test

import (
	"errors"

	"github.com/kkallday/tracker-cli/application"
	"github.com/kkallday/tracker-cli/fakes"
	"github.com/kkallday/tracker-cli/trackerapi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("App", func() {
	var (
		fakeConfigurationLoader *fakes.ConfigurationLoader
		fakeClient              *fakes.Client
		fakeClientProvider      *fakes.ClientProvider
		fakeLogger              *fakes.Logger

		app application.App
	)

	BeforeEach(func() {
		fakeConfigurationLoader = &fakes.ConfigurationLoader{}
		fakeClient = &fakes.Client{}
		fakeClientProvider = &fakes.ClientProvider{}
		fakeLogger = &fakes.Logger{}

		fakeClientProvider.ClientCall.Returns.Client = fakeClient

		app = application.NewApp(fakeLogger, fakeClientProvider, fakeConfigurationLoader)
	})

	Describe("Run", func() {
		It("retrieves project id and token from env vars", func() {
			err := app.Run()
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeConfigurationLoader.LoadCall.CallCount).To(Equal(1))
		})

		It("initializes client with project ID and token", func() {
			fakeConfigurationLoader.LoadCall.Returns.Configuration = application.Configuration{
				ProjectID: 12345,
				Token:     "some-token",
			}

			err := app.Run()
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeClientProvider.ClientCall.CallCount).To(Equal(1))
			Expect(fakeClientProvider.ClientCall.Receives.ProjectID).To(Equal(12345))
			Expect(fakeClientProvider.ClientCall.Receives.Token).To(Equal("some-token"))
		})

		It("retrieves project stories", func() {
			err := app.Run()
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeClient.ProjectStoriesCall.CallCount).To(Equal(1))
		})

		It("writes title and stories to logger", func() {
			fakeClient.ProjectStoriesCall.Returns.Stories = []trackerapi.Story{
				{109832, "feature", "User can do X", 2},
				{201294, "bug", "something is wrong", 0},
				{838312, "chore", "this is a chore", 0},
			}

			err := app.Run()
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeLogger.LogCall.CallCount).To(Equal(1))
			Expect(fakeLogger.LogCall.Receives.Message).To(Equal("Stories in-flight:"))

			Expect(fakeLogger.LogStoriesCall.CallCount).To(Equal(1))

			expectedStories := []trackerapi.Story{
				{109832, "feature", "User can do X", 2},
				{201294, "bug", "something is wrong", 0},
				{838312, "chore", "this is a chore", 0},
			}
			Expect(fakeLogger.LogStoriesCall.Receives.Stories).To(Equal(expectedStories))
		})

		Context("failure cases", func() {
			Context("when configuration loader fails", func() {
				It("returns an error", func() {
					fakeConfigurationLoader.LoadCall.Returns.Error = errors.New("load failed")

					err := app.Run()
					Expect(err).To(MatchError("load failed"))
				})
			})

			Context("when client fails to retrieve project stories", func() {
				It("returns an error", func() {
					fakeClient.ProjectStoriesCall.Returns.Error = errors.New("failed to retrieve project stories")

					err := app.Run()
					Expect(err).To(MatchError("failed to retrieve project stories"))
				})
			})
		})
	})
})
