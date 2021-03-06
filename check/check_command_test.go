package check_test

import (
	"errors"

	"github.com/jamiemonserrate/bintray-resource/bintrayresource"
	"github.com/jamiemonserrate/bintray-resource/check"
	"github.com/jamiemonserrate/bintray-resource/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CheckCommand", func() {
	var (
		fakeBintrayClient *fakes.BintrayClient
		checkRequest      check.CheckRequest
	)

	BeforeEach(func() {
		fakeBintrayClient = &fakes.BintrayClient{
			VersionsToReturn: []string{"1.0.0", "0.0.2", "0.0.1"},
		}
		checkRequest = validCheckRequest()
	})

	It("Requests for the correct package", func() {
		checkRequest.Source.PackageName = "awesome-package"

		checkCommand := check.NewCheckCommand(fakeBintrayClient)
		_, err := checkCommand.Execute(checkRequest)
		Expect(err).ToNot(HaveOccurred())

		Expect(fakeBintrayClient.PackageNameRequested).To(Equal("awesome-package"))
	})

	It("Returns empty array when the latest version is provided", func() {
		checkRequest.RawVersion = bintrayresource.Version{Number: "1.0.0"}

		checkCommand := check.NewCheckCommand(fakeBintrayClient)
		checkResponse, err := checkCommand.Execute(checkRequest)
		Expect(err).ToNot(HaveOccurred())

		Expect(checkResponse).To(BeEmpty())
	})

	It("Returns all versions greater than the one provided", func() {
		checkRequest.RawVersion = bintrayresource.Version{Number: "0.0.1"}

		checkCommand := check.NewCheckCommand(fakeBintrayClient)
		checkResponse, err := checkCommand.Execute(checkRequest)
		Expect(err).ToNot(HaveOccurred())

		Expect(checkResponse).To(Equal(check.CheckResponse{
			bintrayresource.Version{Number: "1.0.0"},
			bintrayresource.Version{Number: "0.0.2"}}))
	})

	It("Returns error from the client", func() {
		fakeBintrayClient.ErrorToBeReturned = errors.New("Some error")

		checkCommand := check.NewCheckCommand(fakeBintrayClient)
		_, err := checkCommand.Execute(checkRequest)

		Expect(err).To(MatchError("Some error"))
	})

	It("Returns error when the request is invalid", func() {
		checkRequest := check.CheckRequest{RawVersion: bintrayresource.Version{Number: ""}}

		checkCommand := check.NewCheckCommand(fakeBintrayClient)
		_, err := checkCommand.Execute(checkRequest)

		Expect(err).To(MatchError("Please specify the SubjectName"))
	})
})
