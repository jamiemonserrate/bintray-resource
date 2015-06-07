package in_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jamiemonserrate/bintray-resource/fakes"
	"github.com/jamiemonserrate/bintray-resource/in"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InCommand", func() {
	var (
		fakeBintrayClient fakes.BintrayClient
		destRootDir       string
		err               error
	)

	BeforeEach(func() {
		fakeBintrayClient = fakes.BintrayClient{}
		destRootDir, err = ioutil.TempDir("", "bintray-resource-integration-test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := os.RemoveAll(destRootDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Downloads the correct package", func() {
		inRequest := in.InRequest{Source: in.Source{PackageName: "awesome-package"},
			RawVersion: in.Version{Number: "1.0.0"}}

		inCommand := in.NewInCommand(&fakeBintrayClient)
		inCommand.Execute(inRequest, "on-the-moon")

		Expect(fakeBintrayClient.PackageNameRequested).To(Equal("awesome-package"))
		Expect(fakeBintrayClient.VersionRequested).To(Equal("1.0.0"))
		Expect(fakeBintrayClient.DestinationDirRequested).To(Equal("on-the-moon"))
	})

	It("Creates the directory", func() {
		destDir := filepath.Join(destRootDir, "i-want-to-be-here")
		Expect(destDir).ToNot(BeAnExistingFile())

		inCommand := in.NewInCommand(&fakeBintrayClient)
		inCommand.Execute(in.InRequest{}, destDir)

		Expect(destDir).To(BeAnExistingFile())
	})

	It("Returns the version of the file in the response", func() {
		inRequest := in.InRequest{RawVersion: in.Version{Number: "0.0.1"}}

		inCommand := in.NewInCommand(&fakeBintrayClient)

		Expect(inCommand.Execute(inRequest, "")).To(Equal(in.InResponse{in.Version{Number: "0.0.1"}}))
	})
})
