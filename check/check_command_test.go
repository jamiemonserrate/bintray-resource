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
			VersionsToReturn: []string{"1.0.0", "0.0.2", "0.0.1"},
		}
	})

	It("Requests for the correct package", func() {
		checkRequest := check.CheckRequest{Source: check.Source{PackageName: "awesome-package"},
			RawVersion: check.Version{Number: "1.0.0"}}

		checkCommand := check.NewCheckCommand(&fakeBintrayClient)
		checkCommand.Execute(checkRequest)

		Expect(fakeBintrayClient.PackageNameRequested).To(Equal("awesome-package"))
	})

	It("Returns empty array when the latest version is provided", func() {
		checkRequest := check.CheckRequest{RawVersion: check.Version{Number: "1.0.0"}}

		checkCommand := check.NewCheckCommand(&fakeBintrayClient)

		Expect(checkCommand.Execute(checkRequest)).To(BeEmpty())
	})

	It("Returns all versions greater than the one provided", func() {
		checkRequest := check.CheckRequest{RawVersion: check.Version{Number: "0.0.1"}}

		checkCommand := check.NewCheckCommand(&fakeBintrayClient)

		Expect(checkCommand.Execute(checkRequest)).To(Equal(check.CheckResponse{
			check.Version{Number: "1.0.0"},
			check.Version{Number: "0.0.2"}}))
	})
})
