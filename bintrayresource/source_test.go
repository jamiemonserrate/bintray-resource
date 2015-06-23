package bintrayresource_test

import (
	"github.com/jamiemonserrate/bintray-resource/bintrayresource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Source", func() {
	Context("When all params provided", func() {
		It("returns true if all required values are passed", func() {
			source := validBintraySource()

			Expect(source.IsValid()).To(BeTrue())
		})
	})

	Context("When all params not provided", func() {
		var source bintrayresource.Source

		BeforeEach(func() {
			source = validBintraySource()
		})

		It("returns false if SubjectName is missing and an error message", func() {
			source.SubjectName = ""

			isValid, errMssg := source.IsValid()

			Expect(isValid).To(BeFalse())
			Expect(errMssg).To(Equal("Please specify the SubjectName"))
		})

		It("returns false if RepoName is missing and an error message", func() {
			source.RepoName = ""

			isValid, errMssg := source.IsValid()

			Expect(isValid).To(BeFalse())
			Expect(errMssg).To(Equal("Please specify the RepoName"))
		})

		It("returns false if PackageName is missing and an error message", func() {
			source.PackageName = ""

			isValid, errMssg := source.IsValid()

			Expect(isValid).To(BeFalse())
			Expect(errMssg).To(Equal("Please specify the PackageName"))
		})

		It("returns false if Username is missing and an error message", func() {
			source.Username = ""

			isValid, errMssg := source.IsValid()

			Expect(isValid).To(BeFalse())
			Expect(errMssg).To(Equal("Please specify the Username"))
		})

		It("returns false if APIKey  is missing and an error message", func() {
			source.APIKey = ""

			isValid, errMssg := source.IsValid()

			Expect(isValid).To(BeFalse())
			Expect(errMssg).To(Equal("Please specify the APIKey"))
		})
	})
})

func validBintraySource() bintrayresource.Source {
	return bintrayresource.Source{SubjectName: "something",
		RepoName: "something", PackageName: "something",
		Username: "something", APIKey: "something"}

}
