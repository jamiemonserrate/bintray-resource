package in

import (
	"os"

	"github.com/hashicorp/go-version"
	"github.com/jamiemonserrate/bintray-resource/bintray"
)

type InCommand struct {
	bintrayClient bintray.BintrayClient
}

func NewInCommand(bintrayClient bintray.BintrayClient) InCommand {
	return InCommand{bintrayClient: bintrayClient}
}

func (inCommand *InCommand) Execute(inRequest InRequest, destinationDir string) InResponse {
	os.MkdirAll(destinationDir, 0755)

	inCommand.bintrayClient.DownloadPackage(inRequest.Source.PackageName, inRequest.RawVersion.Number, destinationDir)

	response := InResponse{Version: inRequest.RawVersion}
	return response
}

func toVersion(versionString string) *version.Version {
	version, _ := version.NewVersion(versionString)
	return version
}
