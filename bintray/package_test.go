package bintray_test

import (
	"github.com/hashicorp/go-version"
	"github.com/jamiemonserrate/bintray-resource/bintray"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Package", func() {
	Context(".VersionsSince", func() {
		It("Returns all versions after the provided version", func() {
			bintrayPackage := bintray.Package{RawVersions: []string{"3.3.3", "2.2.2", "1.1.1"}}

			version1 := newVersion("1.1.1")
			expectedVersion2 := newVersion("2.2.2")
			expectedVersion3 := newVersion("3.3.3")
			Expect(bintrayPackage.VersionsSince(version1)).To(Equal([]*version.Version{expectedVersion3, expectedVersion2}))
		})

		It("Returns all versions in sorted order", func() {
			bintrayPackage := bintray.Package{RawVersions: []string{"4.4.4", "2.2.2", "1.1.1", "3.3.3"}}

			version1 := newVersion("1.1.1")
			expectedVersion2 := newVersion("2.2.2")
			expectedVersion3 := newVersion("3.3.3")
			expectedVersion4 := newVersion("4.4.4")
			Expect(bintrayPackage.VersionsSince(version1)).To(Equal([]*version.Version{expectedVersion4, expectedVersion3, expectedVersion2}))
		})

		It("Returns only latest version if none provided", func() {
			bintrayPackage := bintray.Package{RawVersions: []string{"4.4.4", "2.2.2", "1.1.1", "3.3.3"}, RawLatestVersion: "4.4.4"}

			expectedVersion := newVersion("4.4.4")
			Expect(bintrayPackage.VersionsSince(nil)).To(Equal([]*version.Version{expectedVersion}))
		})
	})
})

func newVersion(stringVersion string) *version.Version {
	version, err := version.NewVersion(stringVersion)
	Expect(err).ToNot(HaveOccurred())
	return version
}
