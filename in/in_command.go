package in

import (
	"errors"
	"os"

	"github.com/hashicorp/go-version"
	"github.com/jamiemonserrate/bintray-resource/bintray"
	"github.com/jamiemonserrate/bintray-resource/bintrayresource"
)

type InCommand struct {
	bintrayClient bintray.BintrayClient
}

func NewInCommand(bintrayClient bintray.BintrayClient) InCommand {
	return InCommand{bintrayClient: bintrayClient}
}

func (inCommand *InCommand) Execute(inRequest InRequest, destinationDir string) (*InResponse, error) {
	if isValid, errMssg := inRequest.IsValid(); !isValid {
		return nil, errors.New(errMssg)
	}

	err := os.MkdirAll(destinationDir, 0755)
	if err != nil {
		return nil, err
	}

	err = inCommand.bintrayClient.DownloadPackage(inRequest.Source.PackageName, inRequest.RawVersion.Number, destinationDir)
	if err != nil {
		return nil, err
	}

	response := &InResponse{Version: inRequest.RawVersion,
		Metadata: []bintrayresource.Metadata{
			bintrayresource.Metadata{Name: "url",
				Value: inCommand.bintrayClient.InPackageURL(inRequest.Source.PackageName,
					inRequest.RawVersion.Number)},
		}}
	return response, nil
}

func toVersion(versionString string) *version.Version {
	version, _ := version.NewVersion(versionString)
	return version
}
