package application_test

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/kkallday/tracker-cli/application"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration Loader", func() {
	Describe("Load", func() {
		It("returns configuration", func() {
			pathToConfigFile, err := writeConfigFile(`{"token": "some-token", "project_id": 54, "api_endpoint_override": "http://www.some-endpoint.com"}`)
			Expect(err).NotTo(HaveOccurred())
			defer os.Remove(pathToConfigFile)

			configLoader := application.NewConfigurationLoader()

			actualConfig, err := configLoader.Load(path.Dir(pathToConfigFile))
			Expect(err).NotTo(HaveOccurred())

			expectedConfig := application.Configuration{
				Token:               "some-token",
				ProjectID:           54,
				APIEndpointOverride: "http://www.some-endpoint.com",
			}
			Expect(actualConfig).To(Equal(expectedConfig))
		})

		Context("failure cases", func() {
			Context("when file cannot be opened", func() {
				It("returns an error", func() {
					configLoader := application.NewConfigurationLoader()

					_, err := configLoader.Load("/non/existent/dir/path")
					Expect(err).To(MatchError(ContainSubstring("no such file or directory")))
				})
			})

			Context("when config json cannot be decoded", func() {
				It("returns an error", func() {
					pathToConfigFile, err := writeConfigFile(`not-valid-json`)
					Expect(err).NotTo(HaveOccurred())
					defer os.Remove(pathToConfigFile)

					configLoader := application.NewConfigurationLoader()
					_, err = configLoader.Load(path.Dir(pathToConfigFile))
					Expect(err).To(MatchError("invalid character 'o' in literal null (expecting 'u')"))
				})
			})

			Context("when config file is missing required values", func() {
				It("returns an error", func() {
					pathToConfigFile, err := writeConfigFile(`{"api_endpoint_override": "http://www.some-endpoint.com"}`)
					Expect(err).NotTo(HaveOccurred())
					defer os.Remove(pathToConfigFile)

					configLoader := application.NewConfigurationLoader()
					_, err = configLoader.Load(path.Dir(pathToConfigFile))
					Expect(err).To(MatchError("Configuration must contain a token and a project ID"))
				})
			})
		})
	})
})

func writeConfigFile(configJSON string) (string, error) {
	configDirPath, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}

	configFile, err := os.Create(filepath.Join(configDirPath, "config.json"))
	defer configFile.Close()
	if err != nil {
		return "", err
	}

	_, err = configFile.WriteString(configJSON)
	if err != nil {
		return "", err
	}

	return configFile.Name(), nil
}
