package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TrackerCLI", func() {
	var (
		testServer *httptest.Server
		session    *gexec.Session
	)

	BeforeEach(func() {
		testServer = startTestServer()
	})

	AfterEach(func() {
		testServer.Close()
	})

	It("prints stories in-flight", func() {
		session = executeTrackerCLI([]string{"--api", testServer.URL}, 12345, "some-tracker-api-token")

		expectedStdout := loadFixture("stories.stdout")
		Expect(session.Out.Contents()).To(Equal([]byte(expectedStdout)))
	})

})

func startTestServer() *httptest.Server {
	fixture := loadFixture("project-stories.json")
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		Expect(path).To(Equal("/services/v5/projects/12345/stories"))

		trackerToken := r.Header.Get("X-TrackerToken")
		Expect(trackerToken).To(Equal("some-tracker-api-token"))
		fmt.Fprint(w, fixture)
	}))

	return testServer
}
