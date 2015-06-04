package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"os/exec"
	"time"
	"bytes"
	"encoding/json"

	"github.com/jamiemonserrate/bintray-resource/check"
)

var _ = Describe("check", func(){
	var (
		checkPath string
		err error
		session *gexec.Session
	)

	BeforeEach(func(){
		checkPath, err = gexec.Build("github.com/jamiemonserrate/bintray-resource/cmd/check")
		Expect(err).NotTo(HaveOccurred())
	})

	It("returns empty array if the version provided is the latest", func(){
		checkRequest := check.CheckRequest{
			Version: check.Version{
				Number: "2.1.1",
			},
		}
		stdin := &bytes.Buffer{}
		err := json.NewEncoder(stdin).Encode(checkRequest)
		Ω(err).ShouldNot(HaveOccurred())
		
		command := exec.Command(checkPath)
		command.Stdin = stdin
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session, 5*time.Second).Should(gexec.Exit(0))
		reader := bytes.NewBuffer(session.Buffer().Contents())

		var response check.CheckResponse
		err = json.NewDecoder(reader).Decode(&response)
		Ω(err).ShouldNot(HaveOccurred())

		Ω(response).Should(BeEmpty())
	})

	It("returns the version when there is a version greater than the input", func(){
		checkRequest := check.CheckRequest{
			Version: check.Version{
				Number: "2.1.0",
			},
		}
		stdin := &bytes.Buffer{}
		err := json.NewEncoder(stdin).Encode(checkRequest)
		Ω(err).ShouldNot(HaveOccurred())
		
		command := exec.Command(checkPath)
		command.Stdin = stdin
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session, 5*time.Second).Should(gexec.Exit(0))
		reader := bytes.NewBuffer(session.Buffer().Contents())

		var response check.CheckResponse
		err = json.NewDecoder(reader).Decode(&response)
		Ω(err).ShouldNot(HaveOccurred())

		Ω(response).Should(Equal(check.CheckResponse{
			{
				Number: "2.1.1",
			},
		}))
	})
})

