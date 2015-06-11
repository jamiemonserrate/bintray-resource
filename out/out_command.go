package out

import (
	"github.com/hashicorp/go-version"
	"github.com/jamiemonserrate/bintray-resource/bintray"
)

type OutCommand struct {
	bintrayClient bintray.BintrayClient
}

func NewOutCommand(bintrayClient bintray.BintrayClient) OutCommand {
	return OutCommand{bintrayClient: bintrayClient}
}

func (outCommand *OutCommand) Execute(outRequest OutRequest) OutResponse {
	// os.MkdirAll(destinationDir, 0755)

	outCommand.bintrayClient.UploadPackage(outRequest.Source.PackageName, outRequest.From, "2.2.4")

	// response := InResponse{Version: inRequest.RawVersion}
	// return response
	return OutResponse{Version: Version{Number: "2.2.4"}}
}

func toVersion(versionString string) *version.Version {
	version, _ := version.NewVersion(versionString)
	return version
}
