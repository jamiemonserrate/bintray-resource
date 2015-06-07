package bintray_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jamiemonserrate/bintray-resource/bintray"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Client", func() {
	var (
		server      *ghttp.Server
		client      *bintray.Client
		destRootDir string
		err         error
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = bintray.NewClient(server.URL(), "subject_name", "repo_name")
		destRootDir, err = ioutil.TempDir("", "bintray-resource-integration-test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		server.Close()
		err := os.RemoveAll(destRootDir)
		Expect(err).ToNot(HaveOccurred())
	})

	Context(".GetPackage", func() {
		It("returns the versions", func() {
			expectedPackage := bintray.Package{
				RawVersions:      []string{"6.6.6", "5.5.5"},
				RawLatestVersion: "6.6.6",
			}

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/packages/subject_name/repo_name/package_name"),
				ghttp.RespondWithJSONEncoded(http.StatusOK, expectedPackage),
			))

			bintrayPackage := client.GetPackage("package_name")

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(bintrayPackage).To(Equal(expectedPackage))
		})
	})

	Context(".DownloadPackage", func() {
		It("returns the versions", func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/subject_name/repo_name/version/package_name"),
				ghttp.RespondWith(http.StatusOK, "the downloaded file content"),
			))

			client.DownloadPackage("package_name", "version", destRootDir)

			Expect(server.ReceivedRequests()).To(HaveLen(1))

			downloadedPackage := filepath.Join(destRootDir, "package_name")
			Expect(downloadedPackage).To(BeAnExistingFile())

			contents, err := ioutil.ReadFile(downloadedPackage)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(contents)).To(Equal("the downloaded file content"))
		})
	})
})
