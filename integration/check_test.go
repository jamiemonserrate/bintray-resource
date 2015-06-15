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

var checkPath string

var _ = Describe("check", func() {
	var (
		err error
	)

	BeforeEach(func() {
		checkPath, err = gexec.Build("github.com/jamiemonserrate/bintray-resource/cmd/check")
		Expect(err).NotTo(HaveOccurred())
	})

	It("returns empty array if the version provided is the latest", func() {
		response := execCheckCommandWith(check.CheckRequest{
			RawVersion: bintrayresource.Version{Number: "2.2.3"},
			Source:     bintrayresource.Source{SubjectName: "jamiemonserrate", RepoName: "jamie-concourse", PackageName: "cf-artifactory"},
		})

		Expect(response).To(BeEmpty())
	})

	It("returns all versions greater than provided version", func() {
		response := execCheckCommandWith(check.CheckRequest{
			RawVersion: bintrayresource.Version{Number: "2.1.0"},
			Source:     bintrayresource.Source{SubjectName: "jamiemonserrate", RepoName: "jamie-concourse", PackageName: "cf-artifactory"},
		})

		Expect(response).To(Equal(check.CheckResponse{{Number: "2.2.3"}, {Number: "2.2.2"}, {Number: "2.1.1"}}))
	})

	It("returns only the latest version if input is empty", func() {
		response := execCheckCommandWith(check.CheckRequest{
			Source: bintrayresource.Source{SubjectName: "jamiemonserrate", RepoName: "jamie-concourse", PackageName: "cf-artifactory"},
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
