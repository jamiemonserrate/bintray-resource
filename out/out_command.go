package out

import (
	"io/ioutil"

	"github.com/hashicorp/go-version"
	"github.com/jamiemonserrate/bintray-resource"
	"github.com/jamiemonserrate/bintray-resource/bintray"
)

type OutCommand struct {
	bintrayClient bintray.BintrayClient
}

func NewOutCommand(bintrayClient bintray.BintrayClient) OutCommand {
	return OutCommand{bintrayClient: bintrayClient}
}

func (outCommand *OutCommand) Execute(outRequest OutRequest) (*OutResponse, error) {
	version, err := ioutil.ReadFile(outRequest.VersionFile)
	if err != nil {
		return nil, err
	}
	err = outCommand.bintrayClient.UploadPackage(outRequest.Source.PackageName,
		outRequest.From, string(version))
	if err != nil {
		return nil, err
	}

	return &OutResponse{Version: bintrayresource.Version{Number: string(version)}}, nil
}

func toVersion(versionString string) *version.Version {
	version, _ := version.NewVersion(versionString)
	return version
}
