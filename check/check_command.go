package check

import "github.com/jamiemonserrate/bintray-resource/bintray"

type Source struct {
	SubjectName string `json:"subject_name"`
	RepoName    string `json:"repo_name"`
	PackageName string `json:"package_name"`
}

type Version struct {
	Number string `json:"number"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
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
	bintrayPackage := checkCommand.bintrayClient.GetPackage(checkRequest.Source.PackageName)

	response := CheckResponse{}
	if checkRequest.Version.Number != bintrayPackage.LatestVersion {
		response = append(response, Version{Number: bintrayPackage.Versions[0]})
	}
	return response
}
