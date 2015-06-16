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

		It("returns false if SubjectName is missing", func() {
			source.SubjectName = ""

			Expect(source.IsValid()).To(BeFalse())
		})

		It("returns false if RepoName is missing", func() {
			source.RepoName = ""

			Expect(source.IsValid()).To(BeFalse())
		})

		It("returns false if PackageName is missing", func() {
			source.PackageName = ""

			Expect(source.IsValid()).To(BeFalse())
		})

		It("returns false if Username is missing", func() {
			source.Username = ""

			Expect(source.IsValid()).To(BeFalse())
		})
		It("returns false if APIKey  is missing", func() {
			source.APIKey = ""

			Expect(source.IsValid()).To(BeFalse())
		})
	})
})

func validBintraySource() bintrayresource.Source {
	return bintrayresource.Source{SubjectName: "something",
		RepoName: "something", PackageName: "something",
		Username: "something", APIKey: "something"}

}
