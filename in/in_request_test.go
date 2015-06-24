package in_test

import (
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

			It("returns false and an erro message if required params missing from Source", func() {
				inRequest.Source.RepoName = ""

				isValid, errMssg := inRequest.IsValid()

				Expect(isValid).To(BeFalse())
				Expect(errMssg).To(Equal("Please specify the RepoName"))
			})

			It("returns false and an error message if Version to download is missing", func() {
				inRequest.RawVersion.Number = ""

				isValid, errMssg := inRequest.IsValid()
				Expect(isValid).To(BeFalse())
				Expect(errMssg).To(Equal("Please specify the Version"))
			})
		})
	})
})
