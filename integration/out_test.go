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
	err      error
	outPath  string
	inputDir string
)

var _ = Describe("out", func() {

	BeforeEach(func() {
		outPath, err = gexec.Build("github.com/jamiemonserrate/bintray-resource/cmd/out")
		Expect(err).NotTo(HaveOccurred())
		inputDir, err = ioutil.TempDir("", "bintray-resource-outtegration-test")
	})

	AfterEach(func() {
		err := os.RemoveAll(inputDir)
		Expect(err).ToNot(HaveOccurred())
	})

	FIt("uploads the file", func() {
		ioutil.WriteFile(filepath.Join(inputDir, "fileToUpload.txt"), []byte("some-content"), 0755)
		response := executeCommandWith(out.OutRequest{
			From:   filepath.Join(inputDir, "fileToUpload.txt"),
			Source: out.Source{SubjectName: "jamiemonserrate", RepoName: "jamie-concourse", PackageName: "cf-artifactory"},
		})

		Expect(response).To(Equal(out.OutResponse{Version: out.Version{Number: "2.2.4"}}))
		downloadDir, err := ioutil.TempDir("", "bintray-resource-outtegration-test-download")
		Expect(err).ToNot(HaveOccurred())
		bintrayClient := bintray.NewClient(
			"https://dl.bintray.com",
			"jamiemonserrate",
			"jamie-concourse")
		bintrayClient.DownloadPackage("cf-artifactory", "2.2.4", downloadDir)
		contents, err := ioutil.ReadFile(filepath.Join(downloadDir, "cf-artifactory"))
		Expect(contents).To(Equal([]byte("some-content")))
		abintrayClient := bintray.NewClient(
			"https://api.bintray.com",
			"jamiemonserrate",
			"jamie-concourse")
		abintrayClient.DeleteVersion("cf-artifactory", "2.2.4")
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
