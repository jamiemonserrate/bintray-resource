package out

import (
	"errors"
	"io/ioutil"

	"github.com/hashicorp/go-version"
	"github.com/jamiemonserrate/bintray-resource/bintray"
	"github.com/jamiemonserrate/bintray-resource/bintrayresource"
)

type OutCommand struct {
	bintrayClient bintray.BintrayClient
}

func NewOutCommand(bintrayClient bintray.BintrayClient) OutCommand {
	return OutCommand{bintrayClient: bintrayClient}
}

func (outCommand *OutCommand) Execute(outRequest OutRequest) (*OutResponse, error) {
	if isValid, errMssg := outRequest.IsValid(); !isValid {
		return nil, errors.New(errMssg)
	}

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
