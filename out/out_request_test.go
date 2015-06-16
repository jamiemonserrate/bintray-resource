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

			It("returns false if required params missing from Source", func() {
				outRequest.Source.RepoName = ""

				Expect(outRequest.IsValid()).To(BeFalse())
			})

			It("returns false if From is not specified", func() {
				outRequest.From = ""

				Expect(outRequest.IsValid()).To(BeFalse())
			})

			It("returns false if Version is not specified", func() {
				outRequest.VersionFile = ""

				Expect(outRequest.IsValid()).To(BeFalse())
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
