package bintray_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/jamiemonserrate/bintray-resource/bintray"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Client", func() {
	var (
		server *ghttp.Server
		client *bintray.Client
		tmpDir string
		err    error
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = bintray.NewClient(server.URL(), "subject_name", "repo_name", "thedude", "topsecretpassword")
		tmpDir, err = ioutil.TempDir("", "bintray-resource-integration-test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		server.Close()
		err := os.RemoveAll(tmpDir)
		Expect(err).ToNot(HaveOccurred())
	})

	Context(".GetPackage", func() {
		It("returns the versions", func() {
			expectedPackage := &bintray.Package{
				RawVersions:      []string{"6.6.6", "5.5.5"},
				RawLatestVersion: "6.6.6",
			}

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/packages/subject_name/repo_name/package_name"),
				ghttp.RespondWithJSONEncoded(http.StatusOK, expectedPackage),
			))

			bintrayPackage, err := client.GetPackage("package_name")

			Expect(err).ToNot(HaveOccurred())
			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(bintrayPackage).To(Equal(expectedPackage))
		})

		Context("handles errors", func() {
			It("if bintray returns error", func() {
				client = bintray.NewClient("some-invalid-url", "doesntmatter", "doesntmatter", "doesntmatter", "doesntmatter")

				bintrayPackage, err := client.GetPackage("package_name")

				Expect(bintrayPackage).To(BeNil())
				Expect(err).To(MatchError(ContainSubstring("unsupported protocol scheme")))
			})

			It("if bintray returns a non 200 status code with a valid error message", func() {
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.RespondWithJSONEncoded(http.StatusNotFound, bintray.ErrorResponse{Message: "The requested path was not found"}),
				))

				bintrayPackage, err := client.GetPackage("package_name")

				Expect(bintrayPackage).To(BeNil())
				Expect(err).To(MatchError("The requested path was not found"))
			})

			It("if bintray returns a 200 status code with a valid error message", func() {
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.RespondWithJSONEncoded(http.StatusOK, bintray.ErrorResponse{Message: "Please try again later"}),
				))

				bintrayPackage, err := client.GetPackage("package_name")

				Expect(bintrayPackage).To(BeNil())
				Expect(err).To(MatchError("Please try again later"))
			})

			Context("when Response cannot be parsed", func() {
				It("returns error for 200 status code", func() {
					server.AppendHandlers(ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, "something we dont expect"),
					))

					bintrayPackage, err := client.GetPackage("package_name")

					Expect(bintrayPackage).To(BeNil())
					Expect(err).To(MatchError("something we dont expect"))
				})

				It("returns error for non 200 status code", func() {
					server.AppendHandlers(ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusInternalServerError, "something we dont expect"),
					))

					bintrayPackage, err := client.GetPackage("package_name")

					Expect(bintrayPackage).To(BeNil())
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError("something we dont expect"))
				})
			})
		})
	})

	Context(".DownloadPackage", func() {
		It("returns the versions", func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/subject_name/repo_name/version/package_name"),
				ghttp.RespondWith(http.StatusOK, "the downloaded file content"),
			))

			err := client.DownloadPackage("package_name", "version", tmpDir)
			Expect(err).ToNot(HaveOccurred())

			Expect(server.ReceivedRequests()).To(HaveLen(1))

			downloadedPackage := filepath.Join(tmpDir, "package_name")
			Expect(downloadedPackage).To(BeAnExistingFile())

			contents, err := ioutil.ReadFile(downloadedPackage)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(contents)).To(Equal("the downloaded file content"))
		})

		Context("handles errors", func() {
			It("returns error if cannot create file", func() {
				err := client.DownloadPackage("doesntmatter", "doesntmatter", "some-invalid-directory")
				Expect(err).To(MatchError(ContainSubstring("no such file or directory")))
			})

			It("if bintray returns error", func() {
				client = bintray.NewClient("some-invalid-url", "doesntmatter", "doesntmatter", "doesntmatter", "doesntmatter")

				err := client.DownloadPackage("doesntmatter", "doesntmatter", tmpDir)

				Expect(err).To(MatchError(ContainSubstring("unsupported protocol scheme")))
			})

			It("if bintray returns a non 200 status code with a valid error message", func() {
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.RespondWithJSONEncoded(http.StatusNotFound, bintray.ErrorResponse{Message: "The requested path was not found"}),
				))

				err := client.DownloadPackage("doesntmatter", "doesntmatter", tmpDir)

				Expect(err).To(MatchError("The requested path was not found"))
			})

			It("if bintray returns a non 200 status code with invalid JSON", func() {
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusNotFound, "The requested path was not found"),
				))

				err := client.DownloadPackage("doesntmatter", "doesntmatter", tmpDir)

				Expect(err).To(MatchError("The requested path was not found"))
			})
		})
	})

	Context(".UploadPackage", func() {
		It("uploads package file to bintray", func() {
			fileToUploadPath := path.Join(tmpDir, "some-package")
			err := ioutil.WriteFile(fileToUploadPath, []byte("super duper package contents"), 0755)
			Expect(err).ToNot(HaveOccurred())
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/content/subject_name/repo_name/package_name/version/version/package_name"),
				VerifyContentsUploaded("super duper package contents"),
				ghttp.VerifyBasicAuth("thedude", "topsecretpassword"),
				ghttp.RespondWith(http.StatusCreated, nil),
			))

			err = client.UploadPackage("package_name", fileToUploadPath, "version")
			Expect(err).ToNot(HaveOccurred())

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})
	})

	Context(".DeleteVersion", func() {
		It("deletes the version from bintray", func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/packages/subject_name/repo_name/package_name/versions/version"),
				ghttp.VerifyBasicAuth("thedude", "topsecretpassword"),
				ghttp.RespondWith(http.StatusOK, nil),
			))

			err = client.DeleteVersion("package_name", "version")
			Expect(err).ToNot(HaveOccurred())

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})
	})
})

func VerifyContentsUploaded(expectedContents string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		Expect(req.Header.Get("Content-Length")).Should(Equal(strconv.Itoa(len(expectedContents))))
		body, err := ioutil.ReadAll(req.Body)
		req.Body.Close()
		Expect(err).ToNot(HaveOccurred())
		Expect(body).ToNot(Equal(expectedContents))
	}
}
