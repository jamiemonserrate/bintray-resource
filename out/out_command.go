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

	ver, err := ioutil.ReadFile(outRequest.VersionFile)
	if err != nil {
		return nil, err
	}

	if _, err := version.NewVersion(string(ver)); err != nil {
		return nil, err
	}

	err = outCommand.bintrayClient.UploadPackage(outRequest.Source.PackageName,
		outRequest.From, string(ver))
	if err != nil {
		return nil, err
	}

	return &OutResponse{Version: bintrayresource.Version{Number: string(ver)}}, nil
}

func toVersion(versionString string) *version.Version {
	version, _ := version.NewVersion(versionString)
	return version
}
