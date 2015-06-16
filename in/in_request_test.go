package in_test

import (
	"github.com/jamiemonserrate/bintray-resource/bintrayresource"
	"github.com/jamiemonserrate/bintray-resource/in"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InRequest", func() {
	Context(".IsValid", func() {
		It("returns true if it has all required params", func() {
			inRequest := validInRequest()
			Expect(inRequest.IsValid()).To(BeTrue())
		})

		Context("when invalid", func() {
			var inRequest in.InRequest

			BeforeEach(func() {
				inRequest = validInRequest()
			})

			It("returns false if required params missing from Source", func() {
				inRequest.Source.RepoName = ""

				Expect(inRequest.IsValid()).To(BeFalse())
			})

			It("returns false if Version to download is missing", func() {
				inRequest.RawVersion.Number = ""

				Expect(inRequest.IsValid()).To(BeFalse())
			})
		})
	})
})

func validInRequest() in.InRequest {
	return in.InRequest{Source: bintrayresource.Source{SubjectName: "something",
		RepoName: "something", PackageName: "something",
		Username: "something", APIKey: "something"},
		RawVersion: bintrayresource.Version{Number: "something"}}
}
