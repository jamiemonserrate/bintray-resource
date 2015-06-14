package integration_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/jamiemonserrate/bintray-resource/bintray"
	"github.com/jamiemonserrate/bintray-resource/out"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	err                error
	outPath            string
	inputDir           string
	versionFilePath    string
	fileToUploadPath   string
	bintrayUsername    = "jamiemonserrate"
	bintrayAPIKey      = "9dd0d7a78b11e773ef4dbc389cf36c1cfe536ebc"
	bintraySubjectName = "jamiemonserrate"
	bintrayRepoName    = "jamie-concourse"
	packageName        = "cf-artifactory"
	expectedVersion    = "2.2.5"
)

var _ = Describe("out", func() {
	var bintrayClient *bintray.Client

	BeforeEach(func() {
		outPath, err = gexec.Build("github.com/jamiemonserrate/bintray-resource/cmd/out")
		Expect(err).NotTo(HaveOccurred())

		inputDir, err = ioutil.TempDir("", "bintray-resource-outtegration-test")
		versionFilePath = filepath.Join(inputDir, "number")
		fileToUploadPath = filepath.Join(inputDir, "fileToUpload.txt")
		bintrayClient = newAPIClient()

		bintrayClient.DeleteVersion(packageName, expectedVersion)
	})

	AfterEach(func() {
		err := os.RemoveAll(inputDir)
		Expect(err).ToNot(HaveOccurred())

		bintrayClient.DeleteVersion(packageName, expectedVersion)
	})

	It("uploads the file and returns version in the response", func() {
		ioutil.WriteFile(fileToUploadPath, []byte("some-content"), 0755)
		ioutil.WriteFile(versionFilePath, []byte(expectedVersion), 0755)

		response := executeCommandWith(out.OutRequest{
			From:        fileToUploadPath,
			VersionFile: versionFilePath,
			Source:      source()})

		Expect(response).To(Equal(out.OutResponse{Version: out.Version{Number: expectedVersion}}))
		Expect(downloadAndReadContentsOf(packageName, expectedVersion)).To(Equal("some-content"))
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

func newAPIClient() *bintray.Client {
	return bintray.NewClient(
		bintray.APIURL,
		bintraySubjectName,
		bintrayRepoName,
		bintrayUsername,
		bintrayAPIKey)
}

func source() out.Source {
	return out.Source{SubjectName: bintraySubjectName,
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
