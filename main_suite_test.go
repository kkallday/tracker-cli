package main_test

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"testing"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTrackercli(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

var pathToTrackerCLI string

var _ = BeforeSuite(func() {
	var err error
	pathToTrackerCLI, err = gexec.Build("github.com/kkallday/tracker-cli")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func executeTrackerCLI(args []string) *gexec.Session {
	command := exec.Command(pathToTrackerCLI, args...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session).Should(gexec.Exit(0))
	return session
}

func loadFixture(pathToFixture string) string {
	fixtureFileContents, err := ioutil.ReadFile(fmt.Sprintf("fixtures/%s", pathToFixture))
	Expect(err).NotTo(HaveOccurred())

	return string(fixtureFileContents)
}
