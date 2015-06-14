package integration_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/jamiemonserrate/bintray-resource"
	"github.com/jamiemonserrate/bintray-resource/in"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	destDir string
	inPath  string
)

var _ = Describe("in", func() {
	var (
		err error
	)

	BeforeEach(func() {
		inPath, err = gexec.Build("github.com/jamiemonserrate/bintray-resource/cmd/in")
		Expect(err).NotTo(HaveOccurred())
		destDir, err = ioutil.TempDir("", "bintray-resource-integration-test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := os.RemoveAll(destDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("Downloads file for the version", func() {
		response := execInCommandWith(in.InRequest{
			RawVersion: bintrayresource.Version{Number: "2.2.3"},
			Source:     bintrayresource.Source{SubjectName: "jamiemonserrate", RepoName: "jamie-concourse", PackageName: "cf-artifactory"},
		})

		Expect(response).To(Equal(in.InResponse{Version: bintrayresource.Version{Number: "2.2.3"}}))

		Expect(filepath.Join(destDir, "cf-artifactory")).To(BeARegularFile())
		contents, err := ioutil.ReadFile(filepath.Join(destDir, "cf-artifactory"))
		Expect(err).ToNot(HaveOccurred())
		Expect(contents).To(Equal([]byte("These contents are valid\n")))
	})
})

func execInCommandWith(inRequest in.InRequest) in.InResponse {
	command := exec.Command(inPath, destDir)
	command.Stdin = encodeInRequest(inRequest)

	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 5*time.Second).Should(gexec.Exit(0))

	return decodeInResponse(session.Buffer().Contents())
}

func encodeInRequest(inRequest in.InRequest) *bytes.Buffer {
	encodedJson := &bytes.Buffer{}
	err := json.NewEncoder(encodedJson).Encode(inRequest)
	Expect(err).ToNot(HaveOccurred())
	return encodedJson
}

func decodeInResponse(encodedResponse []byte) in.InResponse {
	decodedResponse := in.InResponse{}
	err := json.NewDecoder(bytes.NewBuffer(encodedResponse)).Decode(&decodedResponse)
	Expect(err).ToNot(HaveOccurred())
	return decodedResponse
}
