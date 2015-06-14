package check_test

import (
	"errors"

	"github.com/jamiemonserrate/bintray-resource"
	"github.com/jamiemonserrate/bintray-resource/check"
	"github.com/jamiemonserrate/bintray-resource/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CheckCommand", func() {
	var fakeBintrayClient *fakes.BintrayClient

	BeforeEach(func() {
		fakeBintrayClient = &fakes.BintrayClient{
			VersionsToReturn: []string{"1.0.0", "0.0.2", "0.0.1"},
		}
	})

	It("Requests for the correct package", func() {
		checkRequest := check.CheckRequest{Source: bintrayresource.Source{PackageName: "awesome-package"},
			RawVersion: bintrayresource.Version{Number: "1.0.0"}}

		checkCommand := check.NewCheckCommand(fakeBintrayClient)
		_, err := checkCommand.Execute(checkRequest)
		Expect(err).ToNot(HaveOccurred())

		Expect(fakeBintrayClient.PackageNameRequested).To(Equal("awesome-package"))
	})

	It("Returns empty array when the latest version is provided", func() {
		checkRequest := check.CheckRequest{RawVersion: bintrayresource.Version{Number: "1.0.0"}}

		checkCommand := check.NewCheckCommand(fakeBintrayClient)
		checkResponse, err := checkCommand.Execute(checkRequest)
		Expect(err).ToNot(HaveOccurred())

		Expect(checkResponse).To(BeEmpty())
	})

	It("Returns all versions greater than the one provided", func() {
		checkRequest := check.CheckRequest{RawVersion: bintrayresource.Version{Number: "0.0.1"}}

		checkCommand := check.NewCheckCommand(fakeBintrayClient)
		checkResponse, err := checkCommand.Execute(checkRequest)
		Expect(err).ToNot(HaveOccurred())

		Expect(checkResponse).To(Equal(check.CheckResponse{
			bintrayresource.Version{Number: "1.0.0"},
			bintrayresource.Version{Number: "0.0.2"}}))
	})

	It("Returns error from the client", func() {
		checkRequest := check.CheckRequest{RawVersion: bintrayresource.Version{Number: "0.0.1"}}
		fakeBintrayClient.ErrorToBeReturned = errors.New("Some error")

		checkCommand := check.NewCheckCommand(fakeBintrayClient)
		_, err := checkCommand.Execute(checkRequest)

		Expect(err).To(MatchError("Some error"))
	})
})
