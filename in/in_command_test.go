package in_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jamiemonserrate/bintray-resource/bintrayresource"
	"github.com/jamiemonserrate/bintray-resource/fakes"
	"github.com/jamiemonserrate/bintray-resource/in"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InCommand", func() {
	var (
		fakeBintrayClient *fakes.BintrayClient
		destRootDir       string
		err               error
		inRequest         in.InRequest
	)

	BeforeEach(func() {
		fakeBintrayClient = &fakes.BintrayClient{}
		destRootDir, err = ioutil.TempDir("", "bintray-resource-integration-test")
		Expect(err).NotTo(HaveOccurred())

		inRequest = validInRequest()
	})

	AfterEach(func() {
		err := os.RemoveAll(destRootDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Downloads the correct package", func() {
		destDir := filepath.Join(destRootDir, "on-the-moon")
		inRequest.Source.PackageName = "awesome-package"

		inCommand := in.NewInCommand(fakeBintrayClient)
		_, err := inCommand.Execute(inRequest, destDir)

		Expect(err).ToNot(HaveOccurred())
		Expect(fakeBintrayClient.PackageNameRequested).To(Equal("awesome-package"))
		Expect(fakeBintrayClient.VersionRequested).To(Equal("1.0.0"))
		Expect(fakeBintrayClient.DestinationDirRequested).To(Equal(destDir))
	})

	It("Creates the directory", func() {
		destDir := filepath.Join(destRootDir, "i-want-to-be-here")
		Expect(destDir).ToNot(BeAnExistingFile())

		inCommand := in.NewInCommand(fakeBintrayClient)
		_, err := inCommand.Execute(inRequest, destDir)

		Expect(err).ToNot(HaveOccurred())
		Expect(destDir).To(BeAnExistingFile())
	})

	It("Returns the version of the file and metadata in the response", func() {
		destDir := filepath.Join(destRootDir, "i-want-to-be-here")
		inRequest.RawVersion = bintrayresource.Version{Number: "0.0.1"}

		inCommand := in.NewInCommand(fakeBintrayClient)
		inResponse, _ := inCommand.Execute(inRequest, destDir)

		Expect(inResponse.Version.Number).To(Equal("0.0.1"))
		metadata := inResponse.Metadata[0]
		Expect(metadata.Name).To(Equal("url"))
		Expect(metadata.Value).To(Equal("this-is-the-inpackage-url"))
	})

	It("Returns error from the client", func() {
		destDir := filepath.Join(destRootDir, "i-want-to-be-here")
		fakeBintrayClient.ErrorToBeReturned = errors.New("Some error")

		inCommand := in.NewInCommand(fakeBintrayClient)
		_, err := inCommand.Execute(inRequest, destDir)

		Expect(err).To(MatchError("Some error"))
	})

	It("Returns error if cannot create directory", func() {
		destDir := filepath.Join(destRootDir, "i-want-to-be-here")
		os.MkdirAll(destDir, 0400)

		inCommand := in.NewInCommand(fakeBintrayClient)
		_, err := inCommand.Execute(inRequest, filepath.Join(destDir, "some-path"))

		Expect(err).To(MatchError(ContainSubstring("permission denied")))
	})

	It("Returns error if the request is invalid", func() {
		destDir := filepath.Join(destRootDir, "i-want-to-be-here")
		os.MkdirAll(destDir, 0400)
		inRequest.Source.PackageName = ""

		inCommand := in.NewInCommand(fakeBintrayClient)
		_, err := inCommand.Execute(inRequest, filepath.Join(destDir, "some-path"))

		Expect(err).To(MatchError("Please specify the PackageName"))
	})
})
