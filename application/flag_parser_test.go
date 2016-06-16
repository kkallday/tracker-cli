package application_test

import (
	"github.com/kkallday/tracker-cli/application"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flag Parser", func() {
	Describe("Parse", func() {
		It("returns command line args", func() {
			flagParser := application.NewFlagParser()
			actualCommandLineArgs, err := flagParser.Parse([]string{"--config-dir", "/path/to/config"})
			Expect(err).NotTo(HaveOccurred())

			expectedCommandLineArgs := application.CommandLineArgs{ConfigDir: "/path/to/config"}
			Expect(actualCommandLineArgs).To(Equal(expectedCommandLineArgs))
		})

		Context("failure cases", func() {
			Context("when config dir is not given", func() {
				It("returns an error", func() {
					flagParser := application.NewFlagParser()
					_, err := flagParser.Parse([]string{"--some-arg", "some-value"})
					Expect(err).To(MatchError("missing required flag --config-dir"))
				})
			})
		})
	})
})
