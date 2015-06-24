package integration_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jamiemonserrate/bintray-resource/bintray"
	"github.com/jamiemonserrate/bintray-resource/bintrayresource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var (
	checkPath string
	inPath    string
	outPath   string
	err       error

	tmpDir string

	bintrayUsername    = "jamiemonserrate"
	bintrayAPIKey      = "9dd0d7a78b11e773ef4dbc389cf36c1cfe536ebc"
	bintraySubjectName = "jamiemonserrate"
	bintrayRepoName    = "jamie-concourse"
	packageName        = "cf-artifactory"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

func newAPIClient() *bintray.Client {
	return bintray.NewClient(
		bintray.APIURL,
		bintraySubjectName,
		bintrayRepoName,
		bintrayUsername,
		bintrayAPIKey)
}

var _ = BeforeSuite(func() {
	checkPath, err = gexec.Build("github.com/jamiemonserrate/bintray-resource/cmd/check")
	Expect(err).NotTo(HaveOccurred())
	inPath, err = gexec.Build("github.com/jamiemonserrate/bintray-resource/cmd/in")
	Expect(err).NotTo(HaveOccurred())
	outPath, err = gexec.Build("github.com/jamiemonserrate/bintray-resource/cmd/out")
	Expect(err).NotTo(HaveOccurred())

})

var _ = BeforeEach(func() {
	tmpDir, err = ioutil.TempDir("", "bintray-resource-outtegration-test")
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterEach(func() {
	err := os.RemoveAll(tmpDir)
	Expect(err).ToNot(HaveOccurred())
})

func createVersion(version string) {
	fileToUploadPath := filepath.Join(tmpDir, "fileToUpload.txt")
	ioutil.WriteFile(fileToUploadPath, []byte("These contents are valid"), 0755)
	client := newAPIClient()
	err = client.UploadPackage(packageName, fileToUploadPath, version)
	Expect(err).ToNot(HaveOccurred())
}

func deleteVersion(version string) {
	client := newAPIClient()
	err = client.DeleteVersion(packageName, version)
	Expect(err).ToNot(HaveOccurred())
}

func source() bintrayresource.Source {
	return bintrayresource.Source{SubjectName: bintraySubjectName,
		RepoName:    bintrayRepoName,
		PackageName: packageName,
		Username:    bintrayUsername,
		APIKey:      bintrayAPIKey}
}
