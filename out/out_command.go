package out

import (
	"io/ioutil"

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
	version, _ := ioutil.ReadFile(outRequest.VersionFile)
	outCommand.bintrayClient.UploadPackage(outRequest.Source.PackageName,
		outRequest.From, string(version))

	return OutResponse{Version: Version{Number: string(version)}}
}

func toVersion(versionString string) *version.Version {
	version, _ := version.NewVersion(versionString)
	return version
}
