package out_test

import (
	"github.com/jamiemonserrate/bintray-resource/bintrayresource"
	"github.com/jamiemonserrate/bintray-resource/out"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OutRequest", func() {
	Context(".IsValid", func() {
		It("returns true if it has all required params", func() {
			outRequest := validOutRequest()

			Expect(outRequest.IsValid()).To(BeTrue())
		})

		Context("when invalid", func() {
			var outRequest out.OutRequest

			BeforeEach(func() {
				outRequest = validOutRequest()
			})

			It("returns false and an error message if required params missing from Source", func() {
				outRequest.Source.RepoName = ""

				isValid, errMssg := outRequest.IsValid()

				Expect(isValid).To(BeFalse())
				Expect(errMssg).To(Equal("Please specify the RepoName"))
			})

			It("returns false and an error message if From is not specified", func() {
				outRequest.From = ""

				isValid, errMssg := outRequest.IsValid()

				Expect(isValid).To(BeFalse())
				Expect(errMssg).To(Equal("Please specify the From"))
			})

			It("returns false if Version is not specified", func() {
				outRequest.VersionFile = ""

				isValid, errMssg := outRequest.IsValid()

				Expect(isValid).To(BeFalse())
				Expect(errMssg).To(Equal("Please specify the VersionFile"))
			})
		})
	})
})

func validOutRequest() out.OutRequest {
	return out.OutRequest{Source: bintrayresource.Source{SubjectName: "something",
		RepoName: "something", PackageName: "something",
		Username: "something", APIKey: "something"},
		From:        "something",
		VersionFile: "something"}
}
