package integration_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/jamiemonserrate/bintray-resource/bintray"
	"github.com/jamiemonserrate/bintray-resource/bintrayresource"
	"github.com/jamiemonserrate/bintray-resource/out"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("out", func() {
	var (
		versionFilePath  string
		fileToUploadPath string
		expectedVersion  = "2.2.5"
		bintrayClient    *bintray.Client
	)

	BeforeEach(func() {
		versionFilePath = filepath.Join(tmpDir, "number")
		fileToUploadPath = filepath.Join(tmpDir, "fileToUpload.txt")
		bintrayClient = newAPIClient()

		bintrayClient.DeleteVersion(packageName, expectedVersion)
	})

	Context("when no error occurs", func() {
		AfterEach(func() {
			deleteVersion(expectedVersion)
		})

		It("uploads the file and returns version in the response", func() {
			ioutil.WriteFile(fileToUploadPath, []byte("some-content"), 0755)
			ioutil.WriteFile(versionFilePath, []byte(expectedVersion), 0755)

			response := executeCommandWith(out.OutRequest{
				From:        fileToUploadPath,
				VersionFile: versionFilePath,
				Source:      source()})

			Expect(response).To(Equal(out.OutResponse{Version: bintrayresource.Version{Number: expectedVersion}}))
			Expect(downloadAndReadContentsOf(packageName, expectedVersion)).To(Equal("some-content"))
		})
	})

	Context("when an error occurs", func() {
		It("Fails with non zero status code and prints the error", func() {
			outRequest := out.OutRequest{
				Source: bintrayresource.Source{SubjectName: "nonsense"},
			}
			command := exec.Command(outPath)
			command.Stdin = encodeOutRequest(outRequest)

			buffer := gbytes.NewBuffer()
			session, err := gexec.Start(command, GinkgoWriter, buffer)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, 50*time.Second).Should(gexec.Exit(1))

			Eventually(buffer).Should(gbytes.Say(`error runningCommand:`))
		})
	})
})

func executeCommandWith(outRequest out.OutRequest) out.OutResponse {
	command := exec.Command(outPath)
	command.Stdin = encodeOutRequest(outRequest)

	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 50*time.Second).Should(gexec.Exit(0))

	return decodeOutResponse(session.Buffer().Contents())
}

func encodeOutRequest(outRequest out.OutRequest) *bytes.Buffer {
	encodedJson := &bytes.Buffer{}
	err := json.NewEncoder(encodedJson).Encode(outRequest)
	Expect(err).ToNot(HaveOccurred())
	return encodedJson
}

func decodeOutResponse(encodedResponse []byte) out.OutResponse {
	decodedResponse := out.OutResponse{}
	err := json.NewDecoder(bytes.NewBuffer(encodedResponse)).Decode(&decodedResponse)
	Expect(err).ToNot(HaveOccurred())
	return decodedResponse
}

func source() bintrayresource.Source {
	return bintrayresource.Source{SubjectName: bintraySubjectName,
		RepoName:    bintrayRepoName,
		PackageName: packageName,
		Username:    bintrayUsername,
		APIKey:      bintrayAPIKey}
}

func downloadAndReadContentsOf(packageName, version string) string {
	downloadDir, err := ioutil.TempDir("", "bintray-resource-outtegration-test-download")
	Expect(err).ToNot(HaveOccurred())

	bintrayClient := bintray.NewClient(
		bintray.DownloadURL,
		bintraySubjectName,
		bintrayRepoName,
		bintrayUsername,
		bintrayAPIKey)
	bintrayClient.DownloadPackage(packageName, version, downloadDir)
	contents, err := ioutil.ReadFile(filepath.Join(downloadDir, packageName))
	Expect(err).ToNot(HaveOccurred())

	return string(contents)
}
