package main_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TrackerCLI", func() {
	var (
		testServer      *httptest.Server
		pathToConfigDir string
		session         *gexec.Session
	)

	BeforeEach(func() {
		testServer = startTestServer()
		pathToConfigDir = filepath.Dir(writeConfigFile(testServer.URL))
	})

	AfterEach(func() {
		testServer.Close()
		os.Remove(pathToConfigDir)
	})

	It("prints stories in-flight", func() {
		session = executeTrackerCLI([]string{"--config-dir", pathToConfigDir})

		expectedStdout := loadFixture("stories.stdout")
		Expect(session.Out.Contents()).To(Equal([]byte(expectedStdout)))
	})

})

func startTestServer() *httptest.Server {
	fixture := loadFixture("project-stories.json")
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, fixture)
	}))

	return testServer
}

func writeConfigFile(apiEndpointOverride string) string {
	configDirPath, err := ioutil.TempDir("", "")
	Expect(err).NotTo(HaveOccurred())

	configFile, err := os.Create(filepath.Join(configDirPath, "config.json"))
	Expect(err).NotTo(HaveOccurred())
	defer configFile.Close()

	fileContents := fmt.Sprintf(`{"token": "some-token", "project_id": 105, "api_endpoint_override": %q}`, apiEndpointOverride)
	_, err = configFile.WriteString(fileContents)
	Expect(err).NotTo(HaveOccurred())

	return configFile.Name()
}
