package out_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jamiemonserrate/bintray-resource/fakes"
	"github.com/jamiemonserrate/bintray-resource/out"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OutCommand", func() {
	var (
		fakeBintrayClient *fakes.BintrayClient
		tmpDir            string
		err               error
		outRequest        out.OutRequest
	)

	BeforeEach(func() {
		fakeBintrayClient = &fakes.BintrayClient{}
		tmpDir, err = ioutil.TempDir("", "bintray-resource-integration-test")
		Expect(err).NotTo(HaveOccurred())

		outRequest = validOutRequest()
	})

	AfterEach(func() {
		err := os.RemoveAll(tmpDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Uploads the file with the correct version", func() {
		versionFilePath := filepath.Join(tmpDir, "version_file")
		ioutil.WriteFile(versionFilePath, []byte("6.6.6"), 0755)
		outRequest.Source.PackageName = "awesome-package"
		outRequest.VersionFile = versionFilePath
		outRequest.From = "path/to/file/to/be/uploaded"

		outCommand := out.NewOutCommand(fakeBintrayClient)
		outResponse, err := outCommand.Execute(outRequest)
		Expect(err).ToNot(HaveOccurred())

		Expect(fakeBintrayClient.PackageNameRequested).To(Equal("awesome-package"))
		Expect(fakeBintrayClient.FileToBeUploaded).To(Equal("path/to/file/to/be/uploaded"))
		Expect(fakeBintrayClient.VersionToBeUploaded).To(Equal("6.6.6"))
		Expect(outResponse.Version.Number).To(Equal("6.6.6"))
	})

	It("Returns error if cant open file", func() {
		outRequest.VersionFile = "nonsense"
		fakeBintrayClient.ErrorToBeReturned = errors.New("Some error")

		outCommand := out.NewOutCommand(fakeBintrayClient)
		_, err := outCommand.Execute(outRequest)

		Expect(err).To(MatchError(ContainSubstring("no such file or directory")))
	})

	It("Returns error from the client", func() {
		versionFilePath := filepath.Join(tmpDir, "version_file")
		ioutil.WriteFile(versionFilePath, []byte("6.6.6"), 0755)
		outRequest.VersionFile = versionFilePath
		fakeBintrayClient.ErrorToBeReturned = errors.New("Some error")

		outCommand := out.NewOutCommand(fakeBintrayClient)
		_, err := outCommand.Execute(outRequest)

		Expect(err).To(MatchError("Some error"))
	})

	It("Returns error when the request is invalid", func() {
		outRequest.VersionFile = ""

		outCommand := out.NewOutCommand(fakeBintrayClient)
		_, err := outCommand.Execute(outRequest)

		Expect(err).To(MatchError("Please specify the VersionFile"))
	})

})
