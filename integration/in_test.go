package integration_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/jamiemonserrate/bintray-resource/bintrayresource"
	"github.com/jamiemonserrate/bintray-resource/in"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("in", func() {
	var expectedVersion = "2.2.4"

	BeforeEach(func() {
		createVersion(expectedVersion)
	})

	AfterEach(func() {
		deleteVersion(expectedVersion)
	})

	It("Downloads file for the version", func() {
		response := execInCommandWith(in.InRequest{
			RawVersion: bintrayresource.Version{Number: expectedVersion},
			Source:     bintrayresource.Source{SubjectName: bintraySubjectName, RepoName: bintrayRepoName, PackageName: packageName},
		})

		Expect(response).To(Equal(in.InResponse{Version: bintrayresource.Version{Number: expectedVersion}}))

		Expect(filepath.Join(tmpDir, packageName)).To(BeARegularFile())
		contents, err := ioutil.ReadFile(filepath.Join(tmpDir, packageName))
		Expect(err).ToNot(HaveOccurred())
		Expect(contents).To(Equal([]byte("These contents are valid")))
	})

	Context("when an error occurs", func() {
		It("Fails with non zero status code and prints the error", func() {
			inRequest := in.InRequest{
				Source: bintrayresource.Source{SubjectName: "nonsense"},
			}
			command := exec.Command(inPath, tmpDir)
			command.Stdin = encodeInRequest(inRequest)

			buffer := gbytes.NewBuffer()

			session, err := gexec.Start(command, GinkgoWriter, buffer)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, 5*time.Second).Should(gexec.Exit(1))
			Eventually(buffer).Should(gbytes.Say(`error runningCommand:`))
		})
	})
})

func execInCommandWith(inRequest in.InRequest) in.InResponse {
	command := exec.Command(inPath, tmpDir)
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
