package check

import (
	"errors"

	"github.com/hashicorp/go-version"
	"github.com/jamiemonserrate/bintray-resource/bintray"
	"github.com/jamiemonserrate/bintray-resource/bintrayresource"
)

type CheckCommand struct {
	bintrayClient bintray.BintrayClient
}

func NewCheckCommand(bintrayClient bintray.BintrayClient) CheckCommand {
	return CheckCommand{bintrayClient: bintrayClient}
}

func (checkCommand *CheckCommand) Execute(checkRequest CheckRequest) (CheckResponse, error) {
	if isValid, errMssg := checkRequest.IsValid(); !isValid {
		return nil, errors.New(errMssg)
	}

	bintrayPackage, err := checkCommand.bintrayClient.GetPackage(checkRequest.Source.PackageName)

	if err != nil {
		return nil, err
	}

	response := CheckResponse{}
	for _, v := range bintrayPackage.VersionsSince(checkRequest.Version()) {
		response = append(response, bintrayresource.Version{Number: v.String()})
	}
	return response, nil
}

func toVersion(versionString string) *version.Version {
	version, _ := version.NewVersion(versionString)
	return version
}
