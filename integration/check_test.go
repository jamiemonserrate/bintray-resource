package integration_test

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"github.com/jamiemonserrate/bintray-resource"
	"github.com/jamiemonserrate/bintray-resource/check"
)

var _ = Describe("check", func() {
	BeforeEach(func() {
		createVersion("2.2.1")
		createVersion("2.2.2")
		createVersion("2.2.3")
	})

	AfterEach(func() {
		deleteVersion("2.2.3")
		deleteVersion("2.2.2")
		deleteVersion("2.2.1")
	})

	It("returns empty array if the version provided is the latest", func() {
		response := execCheckCommandWith(check.CheckRequest{
			RawVersion: bintrayresource.Version{Number: "2.2.3"},
			Source:     bintrayresource.Source{SubjectName: bintraySubjectName, RepoName: bintrayRepoName, PackageName: packageName},
		})

		Expect(response).To(BeEmpty())
	})

	It("returns all versions greater than provided version", func() {
		response := execCheckCommandWith(check.CheckRequest{
			RawVersion: bintrayresource.Version{Number: "2.1.0"},
			Source:     bintrayresource.Source{SubjectName: bintraySubjectName, RepoName: bintrayRepoName, PackageName: packageName},
		})

		Expect(response).To(Equal(check.CheckResponse{{Number: "2.2.3"}, {Number: "2.2.2"}, {Number: "2.2.1"}}))
	})

	It("returns only the latest version if input is empty", func() {
		response := execCheckCommandWith(check.CheckRequest{
			Source: bintrayresource.Source{SubjectName: bintraySubjectName, RepoName: bintrayRepoName, PackageName: packageName},
		})

		Expect(response).To(Equal(check.CheckResponse{{Number: "2.2.3"}}))
	})

	Context("when an error occurs", func() {
		It("Fails with non zero status code and prints the error", func() {
			checkRequest := check.CheckRequest{
				Source: bintrayresource.Source{SubjectName: "nonsense"},
			}
			command := exec.Command(checkPath)
			command.Stdin = encodeCheckRequest(checkRequest)

			buffer := gbytes.NewBuffer()

			session, err := gexec.Start(command, GinkgoWriter, buffer)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, 5*time.Second).Should(gexec.Exit(1))
			Eventually(buffer).Should(gbytes.Say(`error runningCommand:`))
		})
	})
})

func execCheckCommandWith(checkRequest check.CheckRequest) check.CheckResponse {
	command := exec.Command(checkPath)
	command.Stdin = encodeCheckRequest(checkRequest)

	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 5*time.Second).Should(gexec.Exit(0))

	return decodeCheckResponse(session.Buffer().Contents())
}

func encodeCheckRequest(checkRequest check.CheckRequest) *bytes.Buffer {
	encodedJson := &bytes.Buffer{}
	err := json.NewEncoder(encodedJson).Encode(checkRequest)
	Expect(err).ToNot(HaveOccurred())
	return encodedJson
}

func decodeCheckResponse(encodedResponse []byte) check.CheckResponse {
	decodedResponse := check.CheckResponse{}
	err := json.NewDecoder(bytes.NewBuffer(encodedResponse)).Decode(&decodedResponse)
	Expect(err).ToNot(HaveOccurred())
	return decodedResponse
}
