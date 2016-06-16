package trackerapi_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/kkallday/tracker-cli/trackerapi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	Describe("ProjectStories", func() {
		It("returns stories", func() {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				storiesJSON, err := storiesJSONFixture()
				Expect(err).NotTo(HaveOccurred())

				fmt.Fprintf(w, storiesJSON)
			}))
			defer testServer.Close()

			client := trackerapi.NewClient(testServer.URL, "")
			actualStories, err := client.ProjectStories(1)
			Expect(err).NotTo(HaveOccurred())

			expectedStories := []trackerapi.Story{
				{Id: 1091909, Story_type: "feature", Name: "User can signup", Estimate: 1},
				{Id: 1909283, Story_type: "feature", Name: "User can create a todo list", Estimate: 2},
				{Id: 1032183, Story_type: "chore", Name: "Refactor app startup script", Estimate: 0},
				{Id: 2308423, Story_type: "bug", Name: "Signup success message should go away after navigating away", Estimate: 0},
			}

			Expect(actualStories).To(Equal(expectedStories))
		})

		It("makes a request to correct path with query and header", func() {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				path := r.URL.Path
				Expect(path).To(Equal("/services/v5/projects/6/stories"))

				query := r.URL.Query().Get("with_state")
				Expect(query).To(Equal("started"))

				trackerToken := r.Header.Get("X-TrackerToken")
				Expect(trackerToken).To(Equal("some-tracker-api-token"))

				fmt.Fprint(w, "[]")
			}))
			defer testServer.Close()

			client := trackerapi.NewClient(testServer.URL, "some-tracker-api-token")
			client.ProjectStories(6)
		})

		Context("failure cases", func() {
			Context("when request creation fails", func() {
				It("returns an error", func() {
					client := trackerapi.NewClient("http://%%%%%", "")

					_, err := client.ProjectStories(9)

					Expect(err).To(MatchError("parse http://%%%%%/services/v5/projects/9/stories?with_state=started: invalid URL escape \"%%%\""))
				})
			})

			Context("when request fails", func() {
				It("returns an error", func() {
					client := trackerapi.NewClient("UNKNOWN-PROTOCOL://foo.com", "")

					_, err := client.ProjectStories(4)

					Expect(err).To(MatchError("Get unknown-protocol://foo.com/services/v5/projects/4/stories?with_state=started: unsupported protocol scheme \"unknown-protocol\""))
				})
			})

			Context("when json decoding fails", func() {
				It("returns an error", func() {
					testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprintf(w, "not-valid-json")
					}))
					defer testServer.Close()

					client := trackerapi.NewClient(testServer.URL, "")

					_, err := client.ProjectStories(1)

					Expect(err).To(MatchError("invalid character 'o' in literal null (expecting 'u')"))
				})
			})

			Context("when response code is not 200 OK", func() {
				It("returns an error", func() {
					testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusNotFound)
						w.Write([]byte(`{"message": "an error occurred"}`))
					}))
					defer testServer.Close()

					client := trackerapi.NewClient(testServer.URL, "")

					_, err := client.ProjectStories(1)

					Expect(err).To(MatchError(`bad response: {"message": "an error occurred"}`))
				})
			})
		})
	})
})

func storiesJSONFixture() (string, error) {
	fixtureFileContents, err := ioutil.ReadFile("fixtures/stories.json")
	if err != nil {
		return "", err
	}

	return string(fixtureFileContents), nil
}
