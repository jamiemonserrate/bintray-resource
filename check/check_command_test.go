package check_test

import (
	"github.com/jamiemonserrate/bintray-resource/check"
	"github.com/jamiemonserrate/bintray-resource/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CheckCommand", func() {
	var fakeBintrayClient fakes.BintrayClient

	BeforeEach(func() {
		fakeBintrayClient = fakes.BintrayClient{
			LatestVersionToReturn: "latest_version",
			VersionsToReturn:      []string{"latest_version", "previous_version"},
		}
	})

	It("Requests for the correct package", func() {

		checkRequest := check.CheckRequest{Source: check.Source{PackageName: "awesome-package"}}

		checkCommand := check.NewCheckCommand(&fakeBintrayClient)
		checkCommand.Execute(checkRequest)

		Expect(fakeBintrayClient.PackageNameRequested).To(Equal("awesome-package"))
	})

	It("Returns empty array when the latest version is provided", func() {
		checkRequest := check.CheckRequest{Version: check.Version{Number: "latest_version"}}

		checkCommand := check.NewCheckCommand(&fakeBintrayClient)

		Expect(checkCommand.Execute(checkRequest)).To(BeEmpty())
	})

	It("Returns the version greater than the one provided", func() {
		checkRequest := check.CheckRequest{Version: check.Version{Number: "previous_version"}}

		checkCommand := check.NewCheckCommand(&fakeBintrayClient)

		Expect(checkCommand.Execute(checkRequest)).To(Equal(check.CheckResponse{check.Version{Number: "latest_version"}}))
	})
})
