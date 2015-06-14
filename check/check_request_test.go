package check_test

import (
	"github.com/hashicorp/go-version"
	"github.com/jamiemonserrate/bintray-resource"
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

})
