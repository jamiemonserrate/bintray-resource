package integration_test

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

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
			Version: check.Version{Number: "2.1.1"},
			Source:  check.Source{SubjectName: "jamiemonserrate", RepoName: "jamie-concourse", PackageName: "cf-artifactory"},
		})

		Expect(response).To(BeEmpty())
	})

	It("returns the version when there is a version greater than the input", func() {
		response := execCheckCommandWith(check.CheckRequest{
			Version: check.Version{Number: "2.1.0"},
			Source:  check.Source{SubjectName: "jamiemonserrate", RepoName: "jamie-concourse", PackageName: "cf-artifactory"},
		})

		Expect(response).To(Equal(check.CheckResponse{{Number: "2.1.1"}}))
	})
})

func execCheckCommandWith(checkRequest check.CheckRequest) check.CheckResponse {
	command := exec.Command(checkPath)
	command.Stdin = encodeJson(checkRequest)

	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 5*time.Second).Should(gexec.Exit(0))

	return decodeJson(session.Buffer().Contents())
}

func encodeJson(checkRequest check.CheckRequest) *bytes.Buffer {
	encodedJson := &bytes.Buffer{}
	err := json.NewEncoder(encodedJson).Encode(checkRequest)
	Expect(err).ToNot(HaveOccurred())
	return encodedJson
}

func decodeJson(encodedResponse []byte) check.CheckResponse {
	decodedResponse := check.CheckResponse{}
	err := json.NewDecoder(bytes.NewBuffer(encodedResponse)).Decode(&decodedResponse)
	Expect(err).ToNot(HaveOccurred())
	return decodedResponse
}
