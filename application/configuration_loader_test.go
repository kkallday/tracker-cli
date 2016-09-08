package application_test

import (
	"errors"

	"github.com/kkallday/tracker-cli/application"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration Loader", func() {
	var configLoader application.ConfigurationLoader
	BeforeEach(func() {
		configLoader = application.NewConfigurationLoader()
	})

	AfterEach(func() {
		application.ResetGetenv()
	})

	Describe("Load", func() {
		It("returns a configuration with values from env vars", func() {
			application.SetGetenv(func(key string) string {
				switch key {
				case "PROJECT_ID":
					return "12345"
				case "TOKEN":
					return "some-token"
				default:
					panic(errors.New("unknown key requested"))
				}
			})

			actualConfig, err := configLoader.Load()
			Expect(err).ToNot(HaveOccurred())

			expectedConfig := application.Configuration{
				ProjectID: 12345,
				Token:     "some-token",
			}
			Expect(actualConfig).To(Equal(expectedConfig))
		})

		Context("failure cases", func() {
			Context("when project ID is not an integer", func() {
				It("returns an error", func() {
					application.SetGetenv(func(key string) string {
						switch key {
						case "PROJECT_ID":
							return "not-a-number"
						default:
							panic(errors.New("unknown key requested"))
						}
					})
					_, err := configLoader.Load()
					Expect(err).To(MatchError(`"not-a-number" is not a valid PROJECT_ID. A number is required.`))
				})
			})

			Context("when project ID is not provided", func() {
				It("returns an error", func() {
					application.SetGetenv(func(key string) string {
						switch key {
						case "PROJECT_ID":
							return ""
						default:
							panic(errors.New("unknown key requested"))
						}
					})

					_, err := configLoader.Load()
					Expect(err).To(MatchError("PROJECT_ID is required."))
				})
			})
		})
	})
})
