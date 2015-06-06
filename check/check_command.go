package check

import (
	"github.com/hashicorp/go-version"
	"github.com/jamiemonserrate/bintray-resource/bintray"
)

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
	if toVersion(checkRequest.Version.Number).LessThan(toVersion(bintrayPackage.LatestVersion)) {
		for _, bintrayVersion := range bintrayPackage.Versions {
			v, _ := version.NewVersion(bintrayVersion)
			if v.GreaterThan(toVersion(checkRequest.Version.Number)) {
				response = append(response, Version{Number: bintrayVersion})
			}
		}
	}
	return response
}

func toVersion(versionString string) *version.Version {
	version, _ := version.NewVersion(versionString)
	return version
}
