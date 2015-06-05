package check

import "github.com/jamiemonserrate/bintray-resource/bintray"

type Version struct {
	Number string `json:"number"`
}

type CheckRequest struct {
	Version Version `json:"version"`
}

type CheckResponse []Version

type CheckCommand struct {
	bintrayClient bintray.BintrayClient
}

func NewCheckCommand(bintrayClient bintray.BintrayClient) CheckCommand {
	return CheckCommand{bintrayClient: bintrayClient}
}

func (checkCommand *CheckCommand) Execute(checkRequest CheckRequest) CheckResponse {
	bintrayPackage := checkCommand.bintrayClient.GetPackage("cf-artifactory")

	response := CheckResponse{}
	if checkRequest.Version.Number != bintrayPackage.LatestVersion {
		response = append(response, Version{Number: bintrayPackage.Versions[0]})
	}
	return response
}
