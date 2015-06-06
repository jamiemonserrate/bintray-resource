package bintray_test

import (
	"net/http"

	"github.com/jamiemonserrate/bintray-resource/bintray"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Client", func() {
	var (
		server *ghttp.Server
		client *bintray.Client
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = bintray.NewClient(server.URL(), "subject_name", "repo_name")
	})

	AfterEach(func() {
		server.Close()
	})

	Context(".GetVersions", func() {
		It("returns the versions", func() {
			expectedPackage := bintray.Package{
				LatestVersion: "6.6.6",
				Versions:      []string{"6.6.6", "5.5.5"},
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
})
