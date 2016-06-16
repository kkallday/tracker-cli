package application_test

import (
	"bytes"

	"github.com/kkallday/tracker-cli/application"
	"github.com/kkallday/tracker-cli/trackerapi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logger", func() {
	Describe("LogStories", func() {
		It("prints stories", func() {
			fakeWriter := bytes.NewBuffer([]byte{})

			logger := application.NewLogger(fakeWriter)

			logger.LogStories(
				trackerapi.Story{109832, "feature", "User can signup", 2},
				trackerapi.Story{201294, "bug", "something is wrong", 0},
			)

			actualContents, err := fakeWriter.ReadBytes(byte('\n'))
			Expect(err).NotTo(HaveOccurred())

			expectedContents := []byte("Feature   #109832   Pts: 2 User can signup\n")
			Expect(actualContents).To(Equal(expectedContents))

			actualContents, err = fakeWriter.ReadBytes(byte('\n'))
			Expect(err).NotTo(HaveOccurred())

			expectedContents = []byte("Bug       #201294   something is wrong\n")
			Expect(actualContents).To(Equal(expectedContents))
		})
	})

	Describe("Log", func() {
		It("prints message", func() {
			fakeWriter := bytes.NewBuffer([]byte{})

			logger := application.NewLogger(fakeWriter)
			logger.Log("some message")

			actualContents, err := fakeWriter.ReadBytes(byte('\n'))
			Expect(err).NotTo(HaveOccurred())

			expectedContents := []byte("some message\n")
			Expect(actualContents).To(Equal(expectedContents))
		})
	})
})
