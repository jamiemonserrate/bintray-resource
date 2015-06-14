package out_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jamiemonserrate/bintray-resource"
	"github.com/jamiemonserrate/bintray-resource/fakes"
	"github.com/jamiemonserrate/bintray-resource/out"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OutCommand", func() {
	var (
		fakeBintrayClient fakes.BintrayClient
		tmpDir            string
		err               error
	)

	BeforeEach(func() {
		fakeBintrayClient = fakes.BintrayClient{}
		tmpDir, err = ioutil.TempDir("", "bintray-resource-integration-test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := os.RemoveAll(tmpDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Uploads the file with the correct version", func() {
		versionFilePath := filepath.Join(tmpDir, "version_file")
		ioutil.WriteFile(versionFilePath, []byte("6.6.6"), 0755)
		outRequest := out.OutRequest{Source: bintrayresource.Source{PackageName: "awesome-package"},
			VersionFile: versionFilePath,
			From:        "path/to/file/to/be/uploaded"}

		outCommand := out.NewOutCommand(&fakeBintrayClient)
		outResponse := outCommand.Execute(outRequest)

		Expect(fakeBintrayClient.PackageNameRequested).To(Equal("awesome-package"))
		Expect(fakeBintrayClient.FileToBeUploaded).To(Equal("path/to/file/to/be/uploaded"))
		Expect(fakeBintrayClient.VersionToBeUploaded).To(Equal("6.6.6"))
		Expect(outResponse.Version.Number).To(Equal("6.6.6"))
	})

})
