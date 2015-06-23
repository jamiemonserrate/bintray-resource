package check_test

import (
	"github.com/hashicorp/go-version"
	"github.com/jamiemonserrate/bintray-resource/bintrayresource"
	"github.com/jamiemonserrate/bintray-resource/check"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CheckRequest", func() {
	Context(".Version", func() {
		It("returns version as an object", func() {
			checkRequest := check.CheckRequest{RawVersion: bintrayresource.Version{Number: "1.1.1"}}

			expectedVersion, _ := version.NewVersion("1.1.1")
			Expect(checkRequest.Version()).To(Equal(expectedVersion))
		})
	})

	Context(".IsValid", func() {
		It("returns true if it has all required params", func() {
			checkRequest := validCheckRequest()

			Expect(checkRequest.IsValid()).To(BeTrue())
		})

		It("returns false and an error message if required params missing", func() {
			checkRequest := validCheckRequest()
			checkRequest.Source.RepoName = ""
			isValid, errMessage := checkRequest.IsValid()

			Expect(isValid).To(BeFalse())
			Expect(errMessage).To(Equal("Please specify the RepoName"))
		})
	})
})

func validCheckRequest() check.CheckRequest {
	return check.CheckRequest{Source: bintrayresource.Source{SubjectName: "something",
		RepoName: "something", PackageName: "something",
		Username: "something", APIKey: "something"}}
}
